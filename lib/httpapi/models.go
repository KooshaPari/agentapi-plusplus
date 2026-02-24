package httpapi

import (
	"time"

	mf "github.com/coder/agentapi/lib/msgfmt"
	st "github.com/coder/agentapi/lib/screentracker"
	"github.com/coder/agentapi/lib/util"
	"github.com/danielgtaylor/huma/v2"
)

type MessageType string

const (
	MessageTypeUser    MessageType = "user"
	MessageTypeRaw    MessageType = "raw"
	MessageTypeCommand MessageType = "command"
)

var MessageTypeValues = []MessageType{
	MessageTypeUser,
	MessageTypeRaw,
	MessageTypeCommand,
}

func (m MessageType) Schema(r huma.Registry) *huma.Schema {
	return util.OpenAPISchema(r, "MessageType", MessageTypeValues)
}

// Message represents a message
type Message struct {
	Id      int                 `json:"id" doc:"Unique identifier for the message. This identifier also represents the order of the message in the conversation history."`
	Content string              `json:"content" example:"Hello world" doc:"Message content. The message is formatted as it appears in the agent's terminal session, meaning that, by default, it consists of lines of text with 80 characters per line."`
	Role    st.ConversationRole `json:"role" doc:"Role of the message author"`
	Time    time.Time           `json:"time" doc:"Timestamp of the message"`
}

<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
// LogsResponse represents server logs
type LogsResponse struct {
	Body struct {
		Logs []string `json:"logs" doc:"Server logs"`
	}
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Body struct {
		Status string `json:"status" doc:"Health status"`
	}
}

// ConfigResponse represents the server configuration
type ConfigResponse struct {
=======
// ReadyResponse represents the readiness check response
type ReadyResponse struct {
>>>>>>> 672f6c4 (feat: add /ready endpoint for Kubernetes (#16))
=======
// RateLimitResponse represents rate limit status
type RateLimitResponse struct {
>>>>>>> b81e023 (feat: add rate limiting endpoint (#19))
=======
// RateLimitResponse represents rate limit status
type RateLimitResponse struct {
>>>>>>> d6b99b8 (feat: add rate limiting endpoint)
=======
// LogsResponse represents server logs
type LogsResponse struct {
>>>>>>> cacd4d3 (feat: add logging endpoint (#20))
	Body struct {
		Logs []string `json:"logs" doc:"Server logs"`
	}
}

// StatusResponse represents the server status
type StatusResponse struct {
	Body struct {
		Status    AgentStatus  `json:"status" doc:"Current agent status. 'running' means that the agent is processing a message, 'stable' means that the agent is idle and waiting for input."`
		AgentType mf.AgentType `json:"agent_type" doc:"Type of the agent being used by the server."`
	}
}

// MessagesResponse represents the list of messages
type MessagesResponse struct {
	Body struct {
		Messages []Message `json:"messages" nullable:"false" doc:"List of messages"`
	}
}

<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> 1bc5ee9 (feat: add DELETE /messages to clear history (#15))
=======
>>>>>>> 3372374 (feat: add DELETE /messages to clear conversation history)
// MessagesClearResponse represents the response after clearing messages
type MessagesClearResponse struct {
	Body struct {
		Ok    bool `json:"ok" doc:"Whether messages were cleared"`
		Count int  `json:"count" doc:"Number of messages cleared"`
	}
}

<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> c104df8 (feat: add /messages/count endpoint (#14))
// MessagesCountResponse represents the message count
type MessagesCountResponse struct {
	Body struct {
		Count int `json:"count" doc:"Total number of messages"`
	}
}

=======
>>>>>>> 672f6c4 (feat: add /ready endpoint for Kubernetes (#16))
=======
>>>>>>> 1bc5ee9 (feat: add DELETE /messages to clear history (#15))
=======
>>>>>>> 3372374 (feat: add DELETE /messages to clear conversation history)
type MessageRequestBody struct {
	Content string      `json:"content" example:"/help" doc:"Message content"`
	Type    MessageType `json:"type" doc:"A 'user' type message will be logged as a user message in the conversation history and submitted to the agent. AgentAPI will wait until the agent starts carrying out the task described in the message before responding. A 'raw' type message will be written directly to the agent's terminal session as keystrokes and will not be saved in the conversation history. 'raw' messages are useful for sending escape sequences to the terminal. A 'command' type message sends a slash command directly to the agent (e.g., /help, /resume, /undo)."`
}

// MessageRequest represents a request to create a new message
type MessageRequest struct {
	Body MessageRequestBody `json:"body" doc:"Message content and type"`
}

// MessageResponse represents a newly created message
type MessageResponse struct {
	Body struct {
		Ok bool `json:"ok" doc:"Indicates whether the message was sent successfully. For messages of type 'user', success means detecting that the agent began executing the task described. For messages of type 'raw', success means the keystrokes were sent to the terminal."`
	}
}

type UploadResponse struct {
	Body struct {
		Ok       bool   `json:"ok" doc:"Indicates whether the files were uploaded successfully."`
		FilePath string `json:"filePath" doc:"Path of the file"`
	}
}

type UploadRequest struct {
	File huma.FormFile `form:"file" required:"true" doc:"file that needs to be uploaded"`
}
