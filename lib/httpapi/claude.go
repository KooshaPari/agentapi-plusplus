package httpapi

import (
	mf "github.com/coder/agentapi/lib/msgfmt"
	st "github.com/coder/agentapi/lib/screentracker"
)

func formatPaste(message string) []st.MessagePart {
	return []st.MessagePart{
		// Bracketed paste mode start sequence
		st.MessagePartText{Content: "\x1b[200~", Hidden: true},
		st.MessagePartText{Content: message},
		// Bracketed paste mode end sequence
		st.MessagePartText{Content: "\x1b[201~", Hidden: true},
	}
}

func formatClaudeCodeMessage(message string) []st.MessagePart {
	parts := make([]st.MessagePart, 0)
	// janky hack: send a random character and then a backspace because otherwise
	// Claude Code echoes the startSeq back to the terminal.
	// This basically simulates a user typing and then removing the character.
	parts = append(parts, st.MessagePartText{Content: "x\b", Hidden: true})
	parts = append(parts, formatPaste(message)...)

	return parts
}

func formatPlainMessage(message string) []st.MessagePart {
	return []st.MessagePart{st.MessagePartText{Content: message}}
}

func FormatMessage(agentType mf.AgentType, message string) []st.MessagePart {
	message = mf.TrimWhitespace(message)
	if agentType == mf.AgentTypeCustom {
		return formatPlainMessage(message)
	}
	// For the interactive CLIs we drive via PTY, preserve the existing
	// bracketed-paste flow that avoids partial terminal echo.
	return formatClaudeCodeMessage(message)
}
