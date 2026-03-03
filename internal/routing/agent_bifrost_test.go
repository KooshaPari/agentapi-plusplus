package routing

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestNewAgentBifrost(t *testing.T) {
	bifrost, err := NewAgentBifrost("http://localhost:8080")
	if err != nil {
		t.Fatalf("NewAgentBifrost failed: %v", err)
	}

	if bifrost == nil {
		t.Fatal("Expected non-nil bifrost instance")
	}

	// Test that bifrost is properly initialized by using its public methods
	rule := bifrost.getRule("test-agent")
	if rule.Agent != "test-agent" {
		t.Errorf("Expected agent to be initialized, got %v", rule)
	}
}

func TestGetRule_DefaultRule(t *testing.T) {
	bifrost, _ := NewAgentBifrost("http://localhost:8080")

	rule := bifrost.getRule("non-existent-agent")

	if rule.Agent != "non-existent-agent" {
		t.Errorf("Expected agent to be non-existent-agent, got %s", rule.Agent)
	}

	if rule.PreferredModel != "claude-3-5-sonnet-20241022" {
		t.Errorf("Expected default preferred model, got %s", rule.PreferredModel)
	}

	if len(rule.FallbackModels) != 2 {
		t.Errorf("Expected 2 fallback models, got %d", len(rule.FallbackModels))
	}

	if rule.MaxRetries != 3 {
		t.Errorf("Expected 3 max retries, got %d", rule.MaxRetries)
	}

	if rule.Timeout != 30 {
		t.Errorf("Expected 30 second timeout, got %d", rule.Timeout)
	}
}

func TestSetRule(t *testing.T) {
	bifrost, _ := NewAgentBifrost("http://localhost:8080")

	customRule := RoutingRule{
		Agent:          "custom-agent",
		PreferredModel: "gpt-4o",
		FallbackModels: []string{"claude-3-5-sonnet-20241022"},
		MaxRetries:     5,
		Timeout:        60,
	}

	bifrost.SetRule(customRule)

	retrievedRule := bifrost.getRule("custom-agent")

	if retrievedRule.PreferredModel != "gpt-4o" {
		t.Errorf("Expected preferred model gpt-4o, got %s", retrievedRule.PreferredModel)
	}

	if retrievedRule.MaxRetries != 5 {
		t.Errorf("Expected 5 max retries, got %d", retrievedRule.MaxRetries)
	}

	if retrievedRule.Timeout != 60 {
		t.Errorf("Expected 60 second timeout, got %d", retrievedRule.Timeout)
	}
}

func TestGetOrCreateSession_NewSession(t *testing.T) {
	bifrost, _ := NewAgentBifrost("http://localhost:8080")

	sessionID := bifrost.getOrCreateSession("test-agent")

	if sessionID == "" {
		t.Fatal("Expected non-empty session ID")
	}

	bifrost.sessionsMut.RLock()
	session, exists := bifrost.sessions[sessionID]
	bifrost.sessionsMut.RUnlock()

	if !exists {
		t.Fatal("Expected session to be stored")
	}

	if session.Agent != "test-agent" {
		t.Errorf("Expected agent test-agent, got %s", session.Agent)
	}

	if session.ID != sessionID {
		t.Errorf("Expected session ID %s, got %s", sessionID, session.ID)
	}
}

func TestGetOrCreateSession_ReusesValidSession(t *testing.T) {
	bifrost, _ := NewAgentBifrost("http://localhost:8080")

	sessionID1 := bifrost.getOrCreateSession("test-agent")
	sessionID2 := bifrost.getOrCreateSession("test-agent")

	if sessionID1 != sessionID2 {
		t.Errorf("Expected same session ID, got %s and %s", sessionID1, sessionID2)
	}
}

func TestGetOrCreateSession_CreatesNewSessionAfterExpiry(t *testing.T) {
	bifrost, _ := NewAgentBifrost("http://localhost:8080")

	// Create first session
	sessionID1 := bifrost.getOrCreateSession("test-agent")

	// Manually expire the session
	bifrost.sessionsMut.Lock()
	session := bifrost.sessions[sessionID1]
	session.Started = time.Now().Add(-2 * time.Hour) // More than 1 hour ago
	bifrost.sessionsMut.Unlock()

	// Create new session (should get a new one)
	sessionID2 := bifrost.getOrCreateSession("test-agent")

	if sessionID1 == sessionID2 {
		t.Errorf("Expected different session IDs after expiry")
	}
}

func TestForwardToCliproxy_Success(t *testing.T) {
	// Mock cliproxy server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(RoutingResponse{
			ID:    "resp_123",
			Model: "claude-3-5-sonnet-20241022",
			Choices: []Choice{{
				Message: Message{
					Role:    "assistant",
					Content: "Hello!",
				},
			}},
			Usage: Usage{
				PromptTokens:     10,
				CompletionTokens: 5,
				TotalTokens:      15,
			},
		})
	}))
	defer server.Close()

	bifrost, _ := NewAgentBifrost(server.URL)

	body := map[string]interface{}{
		"model":   "claude-3-5-sonnet-20241022",
		"prompt":  "Hello",
		"agent":   "test-agent",
	}

	resp, err := bifrost.forwardToCliproxy(context.Background(), body)

	if err != nil {
		t.Fatalf("forwardToCliproxy failed: %v", err)
	}

	if resp == nil {
		t.Fatal("Expected non-nil response")
	}

	if resp.ID != "resp_123" {
		t.Errorf("Expected response ID resp_123, got %s", resp.ID)
	}

	if resp.Model != "claude-3-5-sonnet-20241022" {
		t.Errorf("Expected model claude-3-5-sonnet-20241022, got %s", resp.Model)
	}

	if len(resp.Choices) != 1 {
		t.Errorf("Expected 1 choice, got %d", len(resp.Choices))
	}

	if resp.Choices[0].Message.Content != "Hello!" {
		t.Errorf("Expected message content 'Hello!', got %s", resp.Choices[0].Message.Content)
	}
}

func TestForwardToCliproxy_InvalidResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("not json"))
	}))
	defer server.Close()

	bifrost, _ := NewAgentBifrost(server.URL)

	body := map[string]interface{}{
		"model": "claude-3-5-sonnet-20241022",
	}

	_, err := bifrost.forwardToCliproxy(context.Background(), body)

	if err == nil {
		t.Fatal("Expected error for invalid response JSON")
	}
}

func TestRouteRequest_WithDefaultRule(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(RoutingResponse{
			ID:    "resp_456",
			Model: "claude-3-5-sonnet-20241022",
			Choices: []Choice{{
				Message: Message{
					Role:    "assistant",
					Content: "Routed successfully",
				},
			}},
		})
	}))
	defer server.Close()

	bifrost, _ := NewAgentBifrost(server.URL)

	resp, err := bifrost.RouteRequest(context.Background(), "test-agent", "Test prompt")

	if err != nil {
		t.Fatalf("RouteRequest failed: %v", err)
	}

	if resp == nil {
		t.Fatal("Expected non-nil response")
	}

	if resp.ID != "resp_456" {
		t.Errorf("Expected response ID resp_456, got %s", resp.ID)
	}

	if resp.Choices[0].Message.Content != "Routed successfully" {
		t.Errorf("Expected content 'Routed successfully', got %s", resp.Choices[0].Message.Content)
	}
}

func TestRouteRequest_WithCustomRule(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]interface{}
		json.NewDecoder(r.Body).Decode(&body)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Verify the correct model was used
		json.NewEncoder(w).Encode(RoutingResponse{
			ID:    "resp_custom",
			Model: body["model"].(string),
			Choices: []Choice{{
				Message: Message{
					Role:    "assistant",
					Content: "Custom routed",
				},
			}},
		})
	}))
	defer server.Close()

	bifrost, _ := NewAgentBifrost(server.URL)

	customRule := RoutingRule{
		Agent:          "special-agent",
		PreferredModel: "gpt-4o",
		FallbackModels: []string{"claude-3-5-sonnet-20241022"},
		MaxRetries:     5,
		Timeout:        60,
	}
	bifrost.SetRule(customRule)

	resp, err := bifrost.RouteRequest(context.Background(), "special-agent", "Test prompt")

	if err != nil {
		t.Fatalf("RouteRequest failed: %v", err)
	}

	if resp.Model != "gpt-4o" {
		t.Errorf("Expected model gpt-4o, got %s", resp.Model)
	}
}

func TestRoutingResponse_JSONUnmarshal(t *testing.T) {
	jsonStr := `{
		"id": "test_123",
		"model": "claude-3-5-sonnet-20241022",
		"choices": [
			{
				"message": {
					"role": "assistant",
					"content": "test response"
				}
			}
		],
		"usage": {
			"prompt_tokens": 10,
			"completion_tokens": 20,
			"total_tokens": 30
		}
	}`

	var resp RoutingResponse
	err := json.Unmarshal([]byte(jsonStr), &resp)

	if err != nil {
		t.Fatalf("JSON unmarshal failed: %v", err)
	}

	if resp.ID != "test_123" {
		t.Errorf("Expected ID test_123, got %s", resp.ID)
	}

	if resp.Usage.TotalTokens != 30 {
		t.Errorf("Expected 30 total tokens, got %d", resp.Usage.TotalTokens)
	}
}

func TestRouteRequest_CliproxyChatCompletionsContract(t *testing.T) {
	var seenModel string
	var seenMessages []map[string]interface{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/v1/chat/completions" {
			t.Fatalf("expected /v1/chat/completions, got %s", r.URL.Path)
		}
		if got := r.Header.Get("Content-Type"); !strings.Contains(got, "application/json") {
			t.Fatalf("expected application/json content type, got %s", got)
		}

		var payload map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("failed to decode forwarded request: %v", err)
		}
		model, ok := payload["model"].(string)
		if !ok || model == "" {
			t.Fatalf("forwarded request missing string model")
		}
		rawMessages, ok := payload["messages"].([]interface{})
		if !ok || len(rawMessages) == 0 {
			t.Fatalf("forwarded request missing messages array")
		}

		seenModel = model
		seenMessages = make([]map[string]interface{}, 0, len(rawMessages))
		for _, item := range rawMessages {
			msg, ok := item.(map[string]interface{})
			if !ok {
				t.Fatalf("message item is not an object: %T", item)
			}
			seenMessages = append(seenMessages, msg)
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(RoutingResponse{
			ID:    "chatcmpl_contract_123",
			Model: model,
			Choices: []Choice{{
				Message: Message{
					Role:    "assistant",
					Content: "compatible response",
				},
			}},
			Usage: Usage{
				PromptTokens:     7,
				CompletionTokens: 4,
				TotalTokens:      11,
			},
		})
	}))
	defer server.Close()

	bifrost, _ := NewAgentBifrost(server.URL)
	resp, err := bifrost.RouteRequest(context.Background(), "compat-agent", "ping cliproxy")
	if err != nil {
		t.Fatalf("RouteRequest failed: %v", err)
	}
	if resp.ID == "" || resp.Model == "" {
		t.Fatalf("expected non-empty id/model in response")
	}
	if seenModel == "" || len(seenMessages) != 1 {
		t.Fatalf("expected captured forwarded request model/messages")
	}
	if seenMessages[0]["role"] != "user" {
		t.Fatalf("expected user role in forwarded messages, got %v", seenMessages[0]["role"])
	}
	if seenMessages[0]["content"] != "ping cliproxy" {
		t.Fatalf("unexpected forwarded message content: %v", seenMessages[0]["content"])
	}
}

func TestForwardToCliproxy_RejectsContractInvalidResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"id": "chatcmpl_bad_123",
		})
	}))
	defer server.Close()

	bifrost, _ := NewAgentBifrost(server.URL)
	_, err := bifrost.forwardToCliproxy(context.Background(), map[string]interface{}{
		"model":    "gpt-4o",
		"messages": []map[string]string{{"role": "user", "content": "hello"}},
	})
	if err == nil {
		t.Fatalf("expected contract validation error")
	}
	if !strings.Contains(err.Error(), "contract violation") {
		t.Fatalf("expected contract violation error, got %v", err)
	}
}

func TestAgentSession_JSONMarshal(t *testing.T) {
	now := time.Now()
	session := AgentSession{
		ID:       "sess_123",
		Agent:    "test-agent",
		Started:  now,
		Models:   []string{"model1", "model2"},
		Metadata: map[string]interface{}{"key": "value"},
	}

	jsonBytes, err := json.Marshal(session)
	if err != nil {
		t.Fatalf("JSON marshal failed: %v", err)
	}

	var unmarshaled AgentSession
	err = json.Unmarshal(jsonBytes, &unmarshaled)
	if err != nil {
		t.Fatalf("JSON unmarshal failed: %v", err)
	}

	if unmarshaled.ID != "sess_123" {
		t.Errorf("Expected ID sess_123, got %s", unmarshaled.ID)
	}

	if unmarshaled.Agent != "test-agent" {
		t.Errorf("Expected agent test-agent, got %s", unmarshaled.Agent)
	}
}
