package httpapi

import (
	"testing"

	mf "github.com/coder/agentapi/lib/msgfmt"
	st "github.com/coder/agentapi/lib/screentracker"
	"github.com/stretchr/testify/require"
)

func TestFormatMessage_CustomUsesPlainText(t *testing.T) {
	parts := FormatMessage(mf.AgentTypeCustom, "  hello world  ")
	require.Len(t, parts, 1)

	textPart, ok := parts[0].(st.MessagePartText)
	require.True(t, ok)
	require.Equal(t, "hello world", textPart.Content)
	require.False(t, textPart.Hidden)
}

func TestFormatMessage_ClaudeUsesBracketedPaste(t *testing.T) {
	parts := FormatMessage(mf.AgentTypeClaude, "hello world")
	require.Len(t, parts, 4)

	first, ok := parts[0].(st.MessagePartText)
	require.True(t, ok)
	require.True(t, first.Hidden)
	require.Equal(t, "x\b", first.Content)

	second, ok := parts[1].(st.MessagePartText)
	require.True(t, ok)
	require.True(t, second.Hidden)
	require.Equal(t, "\x1b[200~", second.Content)

	third, ok := parts[2].(st.MessagePartText)
	require.True(t, ok)
	require.Equal(t, "hello world", third.Content)
	require.False(t, third.Hidden)

	fourth, ok := parts[3].(st.MessagePartText)
	require.True(t, ok)
	require.True(t, fourth.Hidden)
	require.Equal(t, "\x1b[201~", fourth.Content)
}
