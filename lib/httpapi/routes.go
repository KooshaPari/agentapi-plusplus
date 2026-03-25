package httpapi

import (
	"net/http"
	"net/url"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/sse"
)

// registerRoutes sets up all API endpoints
func (s *Server) registerRoutes() {
	// GET /logs endpoint
	huma.Get(s.api, "/logs", s.getLogs, func(o *huma.Operation) {
		o.Description = "Returns server logs."
	})

	// GET /rate-limit endpoint
	huma.Get(s.api, "/rate-limit", s.getRateLimit, func(o *huma.Operation) {
		o.Description = "Returns rate limit status."
	})

	// GET /config endpoint
	huma.Get(s.api, "/config", s.getConfig, func(o *huma.Operation) {
		o.Description = "Returns the server configuration."
	})

	// GET /health endpoint - liveness probe for load balancers
	huma.Get(s.api, "/health", s.getHealth, func(o *huma.Operation) {
		o.Description = "Health check endpoint for load balancers."
	})
	// GET /version endpoint
	huma.Get(s.api, "/version", s.getVersion, func(o *huma.Operation) {
		o.Description = "Returns the server version."
	})

	// GET /status endpoint
	huma.Get(s.api, "/status", s.getStatus, func(o *huma.Operation) {
		o.Description = "Returns the current status of the agent."
	})
	// GET /info endpoint - returns agent and server info
	huma.Get(s.api, "/info", s.getInfo, func(o *huma.Operation) {
		o.Description = "Returns information about the server and agent."
	})

	// GET /messages endpoint
	// Query params: after (int) - return messages after this ID, limit (int) - limit results
	huma.Get(s.api, "/messages", s.getMessages, func(o *huma.Operation) {
		o.Description = "Returns a list of messages representing the conversation history with the agent. Supports ?after=<id> and ?limit=<n> query parameters for pagination."
	})

	// DELETE /messages endpoint - clear all messages
	huma.Delete(s.api, "/messages", s.clearMessages, func(o *huma.Operation) {
		o.Description = "Clear all messages from conversation history."
	})
	// GET /messages/count endpoint
	huma.Get(s.api, "/messages/count", s.getMessagesCount, func(o *huma.Operation) {
		o.Description = "Returns the count of messages in the conversation."
	})

	// POST /message endpoint
	huma.Post(s.api, "/message", s.createMessage, func(o *huma.Operation) {
		o.Description = "Send a message to the agent. For messages of type 'user', the agent's status must be 'stable' for the operation to complete successfully. Otherwise, this endpoint will return an error."
	})

	huma.Post(s.api, "/upload", s.uploadFiles, func(o *huma.Operation) {
		o.Description = "Upload files to the specified upload path."
	})

	// GET /events endpoint
	sse.Register(s.api, huma.Operation{
		OperationID: "subscribeEvents",
		Method:      http.MethodGet,
		Path:        "/events",
		Summary:     "Subscribe to events",
		Description: "The events are sent as Server-Sent Events (SSE). Initially, the endpoint returns a list of events needed to reconstruct the current state of the conversation and the agent's status. After that, it only returns events that have occurred since the last event was sent.\n\nNote: When an agent is running, the last message in the conversation history is updated frequently, and the endpoint sends a new message update event each time.",
		Middlewares: []func(huma.Context, func(huma.Context)){sseMiddleware},
	}, map[string]any{
		// Mapping of event type name to Go struct for that event.
		"message_update": MessageUpdateBody{},
		"status_change":  StatusChangeBody{},
	}, s.subscribeEvents)

	sse.Register(s.api, huma.Operation{
		OperationID: "subscribeScreen",
		Method:      http.MethodGet,
		Path:        "/internal/screen",
		Summary:     "Subscribe to screen",
		Hidden:      true,
		Middlewares: []func(huma.Context, func(huma.Context)){sseMiddleware},
	}, map[string]any{
		"screen": ScreenUpdateBody{},
	}, s.subscribeScreen)

	s.router.Handle("/", http.HandlerFunc(s.redirectToChat))

	// Serve static files for the chat interface under /chat
	s.registerStaticFileRoutes()
}

// registerStaticFileRoutes sets up routes for serving static files
func (s *Server) registerStaticFileRoutes() {
	chatHandler := FileServerWithIndexFallback(s.chatBasePath)

	// Mount the file server at /chat
	s.router.Handle("/chat", http.StripPrefix("/chat", chatHandler))
	s.router.Handle("/chat/*", http.StripPrefix("/chat", chatHandler))
}

func (s *Server) redirectToChat(w http.ResponseWriter, r *http.Request) {
	rdir, err := url.JoinPath(s.chatBasePath, "embed")
	if err != nil {
		s.logger.Error("Failed to construct redirect URL", "error", err)
		http.Error(w, "Failed to redirect", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, rdir, http.StatusTemporaryRedirect)
}
