package httpapi

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"slices"
	"sort"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/coder/agentapi/internal/version"
	"github.com/coder/agentapi/lib/logctx"
	mf "github.com/coder/agentapi/lib/msgfmt"
	st "github.com/coder/agentapi/lib/screentracker"
	"github.com/coder/agentapi/lib/termexec"
	"github.com/coder/quartz"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"golang.org/x/xerrors"
)

// Server represents the HTTP server
type Server struct {
	router       chi.Router
	api          huma.API
	port         int
	srv          *http.Server
	mu           sync.RWMutex
	logger       *slog.Logger
	conversation st.Conversation
	agentio      *termexec.Process
	agentType    mf.AgentType
	emitter      *EventEmitter
	chatBasePath string
	tempDir      string
	clock        quartz.Clock
}

func (s *Server) NormalizeSchema(schema any) any {
	switch val := (schema).(type) {
	case *any:
		s.NormalizeSchema(*val)
	case []any:
		for i := range val {
			s.NormalizeSchema(&val[i])
		}
		sort.SliceStable(val, func(i, j int) bool {
			return fmt.Sprintf("%v", val[i]) < fmt.Sprintf("%v", val[j])
		})
	case map[string]any:
		for k := range val {
			valUnderKey := val[k]
			s.NormalizeSchema(&valUnderKey)
			val[k] = valUnderKey
		}
	}
	return schema
}

func (s *Server) GetOpenAPI() string {
	jsonBytes, err := s.api.OpenAPI().Downgrade()
	if err != nil {
		return ""
	}
	// unmarshal the json and pretty print it
	var jsonObj any
	if err := json.Unmarshal(jsonBytes, &jsonObj); err != nil {
		return ""
	}

	// Normalize
	normalized := s.NormalizeSchema(jsonObj)

	prettyJSON, err := json.MarshalIndent(normalized, "", "  ")
	if err != nil {
		return ""
	}
	return string(prettyJSON)
}

// That's about 40 frames per second. It's slightly less
// because the action of taking a snapshot takes time too.
const snapshotInterval = 25 * time.Millisecond

type ServerConfig struct {
	AgentType      mf.AgentType
	Process        *termexec.Process
	Port           int
	ChatBasePath   string
	AllowedHosts   []string
	AllowedOrigins []string
	InitialPrompt  string
	Clock          quartz.Clock
}

// Validate allowed hosts don't contain whitespace, commas, schemes, or ports.
// Viper/Cobra use different separators (space for env vars, comma for flags),
// so these characters likely indicate user error.
func parseAllowedHosts(input []string) ([]string, error) {
	if len(input) == 0 {
		return nil, fmt.Errorf("the list must not be empty")
	}
	if slices.Contains(input, "*") {
		return []string{"*"}, nil
	}
	// First pass: whitespace & comma checks (surface these errors first)
	// Viper/Cobra use different separators (space for env vars, comma for flags),
	// so these characters likely indicate user error.
	for _, item := range input {
		for _, r := range item {
			if unicode.IsSpace(r) {
				return nil, fmt.Errorf("'%s' contains whitespace characters, which are not allowed", item)
			}
		}
		if strings.Contains(item, ",") {
			return nil, fmt.Errorf("'%s' contains comma characters, which are not allowed", item)
		}
	}
	// Second pass: scheme check
	for _, item := range input {
		if strings.Contains(item, "http://") || strings.Contains(item, "https://") {
			return nil, fmt.Errorf("'%s' must not include http:// or https://", item)
		}
	}
	hosts := make([]*url.URL, 0, len(input))
	// Third pass: url parse
	for _, item := range input {
		trimmed := strings.TrimSpace(item)
		u, err := url.Parse("http://" + trimmed)
		if err != nil {
			return nil, fmt.Errorf("'%s' is not a valid host: %w", item, err)
		}
		hosts = append(hosts, u)
	}
	// Fourth pass: port check
	for _, u := range hosts {
		if u.Port() != "" {
			return nil, fmt.Errorf("'%s' must not include a port", u.Host)
		}
	}
	hostStrings := make([]string, 0, len(hosts))
	for _, u := range hosts {
		hostStrings = append(hostStrings, u.Hostname())
	}
	return hostStrings, nil
}

// Validate allowed origins
func parseAllowedOrigins(input []string) ([]string, error) {
	if len(input) == 0 {
		return nil, fmt.Errorf("the list must not be empty")
	}
	if slices.Contains(input, "*") {
		return []string{"*"}, nil
	}
	// Viper/Cobra use different separators (space for env vars, comma for flags),
	// so these characters likely indicate user error.
	for _, item := range input {
		for _, r := range item {
			if unicode.IsSpace(r) {
				return nil, fmt.Errorf("'%s' contains whitespace characters, which are not allowed", item)
			}
		}
		if strings.Contains(item, ",") {
			return nil, fmt.Errorf("'%s' contains comma characters, which are not allowed", item)
		}
	}
	origins := make([]string, 0, len(input))
	for _, item := range input {
		trimmed := strings.TrimSpace(item)
		u, err := url.Parse(trimmed)
		if err != nil {
			return nil, fmt.Errorf("'%s' is not a valid origin: %w", item, err)
		}
		origins = append(origins, fmt.Sprintf("%s://%s", u.Scheme, u.Host))
	}
	return origins, nil
}

// NewServer creates a new server instance
func NewServer(ctx context.Context, config ServerConfig) (*Server, error) {
	router := chi.NewMux()

	logger := logctx.From(ctx)

	if config.Clock == nil {
		config.Clock = quartz.NewReal()
	}

	allowedHosts, err := parseAllowedHosts(config.AllowedHosts)
	if err != nil {
		return nil, xerrors.Errorf("failed to parse allowed hosts: %w", err)
	}
	allowedOrigins, err := parseAllowedOrigins(config.AllowedOrigins)
	if err != nil {
		return nil, xerrors.Errorf("failed to parse allowed origins: %w", err)
	}

	logger.Info(fmt.Sprintf("Allowed hosts: %s", strings.Join(allowedHosts, ", ")))
	logger.Info(fmt.Sprintf("Allowed origins: %s", strings.Join(allowedOrigins, ", ")))

	// Enforce allowed hosts in a custom middleware that ignores the port during matching.
	badHostHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Invalid host header. Allowed hosts: "+strings.Join(allowedHosts, ", "), http.StatusBadRequest)
	})
	router.Use(hostAuthorizationMiddleware(allowedHosts, badHostHandler))

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
	router.Use(corsMiddleware.Handler)

	humaConfig := huma.DefaultConfig("AgentAPI", version.Version)
	humaConfig.Info.Description = "HTTP API for Claude Code, Goose, and Aider.\n\nhttps://github.com/coder/agentapi"
	api := humachi.New(router, humaConfig)
	formatMessage := func(message string, userInput string) string {
		return mf.FormatAgentMessage(config.AgentType, message, userInput)
	}

	isAgentReadyForInitialPrompt := func(message string) bool {
		return mf.IsAgentReadyForInitialPrompt(config.AgentType, message)
	}

	formatToolCall := func(message string) (string, []string) {
		return mf.FormatToolCall(config.AgentType, message)
	}

	emitter := NewEventEmitter(WithAgentType(config.AgentType))

	// Format initial prompt into message parts if provided
	var initialPrompt []st.MessagePart
	if config.InitialPrompt != "" {
		initialPrompt = FormatMessage(config.AgentType, config.InitialPrompt)
	}

	conversation := st.NewPTY(ctx, st.PTYConversationConfig{
		AgentType:             config.AgentType,
		AgentIO:               config.Process,
		Clock:                 config.Clock,
		SnapshotInterval:      snapshotInterval,
		ScreenStabilityLength: 2 * time.Second,
		FormatMessage:         formatMessage,
		ReadyForInitialPrompt: isAgentReadyForInitialPrompt,
		FormatToolCall:        formatToolCall,
		InitialPrompt:         initialPrompt,
		Logger:                logger,
	}, emitter)

	// Create temporary directory for uploads
	tempDir, err := os.MkdirTemp("", "agentapi-uploads-")
	if err != nil {
		return nil, xerrors.Errorf("failed to create temporary directory: %w", err)
	}
	logger.Info("Created temporary directory for uploads", "tempDir", tempDir)

	s := &Server{
		router:       router,
		api:          api,
		port:         config.Port,
		conversation: conversation,
		logger:       logger,
		agentio:      config.Process,
		agentType:    config.AgentType,
		emitter:      emitter,
		chatBasePath: strings.TrimSuffix(config.ChatBasePath, "/"),
		tempDir:      tempDir,
		clock:        config.Clock,
	}

	// Register API routes
	s.registerRoutes()

	// Start the conversation polling loop if we have a process.
	// Process is nil only when --print-openapi is used (no agent runs).
	// The process is already running at this point - termexec.StartProcess()
	// blocks until the PTY is created and the process is active. Agent
	// readiness (waiting for the prompt) is handled asynchronously inside
	// conversation.Start() via ReadyForInitialPrompt.
	if config.Process != nil {
		s.conversation.Start(ctx)
	}

	return s, nil
}

// Handler returns the underlying chi.Router for testing purposes.
func (s *Server) Handler() http.Handler {
	return s.router
}

// Start starts the HTTP server
func (s *Server) Start() error {
	addr := fmt.Sprintf(":%d", s.port)
	s.srv = &http.Server{
		Addr:    addr,
		Handler: s.router,
	}

	return s.srv.ListenAndServe()
}

// Stop gracefully stops the HTTP server
func (s *Server) Stop(ctx context.Context) error {
	// Clean up temporary directory
	s.cleanupTempDir()

	if s.srv != nil {
		return s.srv.Shutdown(ctx)
	}
	return nil
}

// cleanupTempDir removes the temporary directory and all its contents
func (s *Server) cleanupTempDir() {
	if err := os.RemoveAll(s.tempDir); err != nil {
		s.logger.Error("Failed to clean up temporary directory", "tempDir", s.tempDir, "error", err)
	} else {
		s.logger.Info("Cleaned up temporary directory", "tempDir", s.tempDir)
	}
}
