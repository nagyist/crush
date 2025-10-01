package chat

import (
	"context"
	"errors"
	"testing"

	"github.com/charmbracelet/crush/internal/llm/agent"
	"github.com/charmbracelet/crush/internal/message"
	"github.com/charmbracelet/crush/internal/session"
	"github.com/charmbracelet/crush/internal/tui/components/chat"
	"github.com/stretchr/testify/assert"
)

type mockAgentRunner struct {
	idOfSessionRanWithin  string
	ranWithMsgText        string
	ranWithMsgAttachments []message.Attachment
}

func (a *mockAgentRunner) Run(ctx context.Context, sessionID string, content string, attachments ...message.Attachment) (<-chan agent.AgentEvent, error) {
	a.idOfSessionRanWithin = sessionID
	a.ranWithMsgText = content
	a.ranWithMsgAttachments = attachments
	return nil, nil
}

type mockCrushCtx struct {
	currentSession session.Session
	runner         *mockAgentRunner
}

func (m *mockCrushCtx) CoderAgent() (agent.Runner, bool) {
	if m.runner == nil {
		return nil, false
	}
	return m.runner, true
}

func (m *mockCrushCtx) MakeSessionCurrent(id string) error {
	return errors.New("not implemented for these tests")
}

func (m *mockCrushCtx) ResolveCurrentSession() (session.Session, error) {
	return m.currentSession, nil
}

func TestOnSendChatMessage_RunsAgentWithMessageData(t *testing.T) {
	currentSess := session.Session{
		ID: "test-chat-session",
	}
	agentRunner := mockAgentRunner{}
	cctx := mockCrushCtx{
		currentSession: currentSess,
		runner:         &agentRunner,
	}
	onSendChatMessage(&cctx, chat.Msg{
		Text: "Fake message for test",
	})

	assert.Equal(t, currentSess.ID, agentRunner.idOfSessionRanWithin)
	assert.Equal(t, "Fake message for test", agentRunner.ranWithMsgText)
	assert.Nil(t, agentRunner.ranWithMsgAttachments)
}
