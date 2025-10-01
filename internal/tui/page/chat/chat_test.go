package chat

import (
	"context"
	"errors"
	"testing"

	"github.com/charmbracelet/crush/internal/llm/agent"
	"github.com/charmbracelet/crush/internal/message"
	"github.com/charmbracelet/crush/internal/session"
	"github.com/charmbracelet/crush/internal/tui/components/chat"
	"github.com/charmbracelet/crush/internal/tui/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockAgentRunner struct {
	idOfSessionRanWithin  string
	ranWithMsgText        string
	ranWithMsgAttachments []message.Attachment
	errOnRun              error
}

func (a *mockAgentRunner) Run(ctx context.Context, sessionID string, content string, attachments ...message.Attachment) (<-chan agent.AgentEvent, error) {
	a.idOfSessionRanWithin = sessionID
	a.ranWithMsgText = content
	a.ranWithMsgAttachments = attachments
	return nil, a.errOnRun
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
	cmds := onSendChatMessage(&cctx, chat.Msg{
		Text: "Fake message for test",
	})

	require.Len(t, cmds, 1)
	assert.IsType(t, chat.SessionSelectedMsg{}, cmds[0]())
	assert.Equal(t, currentSess.ID, agentRunner.idOfSessionRanWithin)
	assert.Equal(t, "Fake message for test", agentRunner.ranWithMsgText)
	assert.Nil(t, agentRunner.ranWithMsgAttachments)
}

func TestOnSendChatMessage_RunsAgentWithMessageDataButAgentFailsToRun(t *testing.T) {
	currentSess := session.Session{
		ID: "test-chat-session",
	}
	agentRunner := mockAgentRunner{
		errOnRun: errors.New("test err agent failed to run"),
	}
	cctx := mockCrushCtx{
		currentSession: currentSess,
		runner:         &agentRunner,
	}
	cmds := onSendChatMessage(&cctx, chat.Msg{
		Text: "Fake message for test",
	})

	require.Len(t, cmds, 2)
	assert.IsType(t, chat.SessionSelectedMsg{}, cmds[0]())
	reportedErrMsg := cmds[1]()
	assert.IsType(t, util.InfoMsg{}, reportedErrMsg)
	assert.Equal(t, util.InfoMsg{
		Type: 2,
		Msg:  "test err agent failed to run",
	}, reportedErrMsg)
	assert.Equal(t, currentSess.ID, agentRunner.idOfSessionRanWithin)
	assert.Equal(t, "Fake message for test", agentRunner.ranWithMsgText)
	assert.Nil(t, agentRunner.ranWithMsgAttachments)
}
