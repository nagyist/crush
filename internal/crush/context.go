package crush

import (
	"context"
	"fmt"

	"github.com/charmbracelet/crush/internal/llm/agent"
	"github.com/charmbracelet/crush/internal/session"
	"github.com/charmbracelet/crush/internal/tui/components/chat"
)

type StorePromptMsg struct {
	msg chat.Msg
}

type PromptHistory struct {
	SessionID int
}

type Context interface {
	CoderAgent() (agent.Service, bool)
	MakeSessionCurrent(id string) error
	ResolveCurrentSession() (session.Session, error)
}

func NewContext(sessRepo session.Repository) Context {
	return &ccontext{
		sessRepo: sessRepo,
	}
}

type ccontext struct {
	sessRepo       session.Repository
	currentSession *session.Session
	coderAgent     agent.Service
}

func (c *ccontext) ResolveCurrentSession() (session.Session, error) {
	if c.currentSession == nil {
		newSession, err := c.sessRepo.Create(context.Background(), "New Session")
		if err != nil {
			return session.Session{}, fmt.Errorf("failed to create session: %w", err)
		}
		c.currentSession = &newSession
	}
	return *c.currentSession, nil
}

func (c *ccontext) MakeSessionCurrent(id string) error {
	return nil
}

func (c *ccontext) CoderAgent() (agent.Service, bool) {
	return nil, false
}
