package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/coder/agentapi/internal/middleware"
	"github.com/coder/agentapi/internal/routing"
	"github.com/go-chi/chi/v5"
)

// Server represents the agentapi HTTP server
type Server struct {
	port             int
	router           *routing.AgentBifrost
	agentHandler     *AgentHandler
	server           *http.Server
	requestIDHandler *middleware.RequestIDHandler
}

// New creates a new agentapi server
func New(port int, router *routing.AgentBifrost) *Server {
	s := &Server{
		port:   port,
		router: router,
	}
	s.agentHandler = NewAgentHandler(router)
	s.requestIDHandler = middleware.NewRequestIDHandler(30 * time.Second)
	return s
}

// Start starts the HTTP server
func (s *Server) Start() error {
	r := chi.NewRouter()
	if err := middleware.ApplyDefaultStack(r); err != nil {
		return err
	}

	middleware.ReadinessCheckRoute(r)

	// Health check
	r.Get("/health", s.health)

	// Agent lifecycle endpoints
	s.agentHandler.RegisterRoutes(r)

	// Agent routing endpoints
	r.Post("/v1/chat/completions", s.wrapHandler(s.chatCompletions))

	// Management endpoints
	r.Route("/admin", func(r chi.Router) {
		r.Get("/rules", s.wrapHandler(s.listRules))
		r.Post("/rules", s.wrapHandler(s.setRule))
		r.Get("/sessions", s.wrapHandler(s.listSessions))
	})

	// Connect to cliproxy+bifrost
	r.HandleFunc("/proxy/*", s.wrapHandler(s.proxy))

	s.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.port),
		Handler: r,
	}

	return s.server.ListenAndServe()
}

func (s *Server) wrapHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.requestIDHandler.WrapHandler(next).ServeHTTP(w, r)
	}
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown() {
	if s.server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		s.server.Shutdown(ctx)
	}
}

func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (s *Server) chatCompletions(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Agent  string `json:"agent"`
		Model  string `json:"model"`
		Prompt string `json:"prompt"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Use "default" if no agent specified
	agent := req.Agent
	if agent == "" {
		agent = "default"
	}

	resp, err := s.router.RouteRequest(r.Context(), agent, req.Prompt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (s *Server) listRules(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"rules": "configured"})
}

func (s *Server) setRule(w http.ResponseWriter, r *http.Request) {
	var rule routing.RoutingRule
	if err := json.NewDecoder(r.Body).Decode(&rule); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s.router.SetRule(rule)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (s *Server) listSessions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"sessions": "active"})
}

func (s *Server) proxy(w http.ResponseWriter, r *http.Request) {
	// Proxy requests to cliproxy+bifrost
	path := chi.URLParam(r, "*")
	log.Printf("Proxying request to: %s", path)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"proxied": path,
		"method":  r.Method,
	})
}
