package crush

import (
	"context"
	"fmt"

	"github.com/charmbracelet/crush/internal/history"
	"github.com/charmbracelet/crush/internal/llm/agent"
	"github.com/charmbracelet/crush/internal/session"
)

type Context interface {
	CoderAgent() (agent.Runner, bool)
	MakeSessionCurrent(id string) error
	ResolveCurrentSession() (session.Session, error)
}

func NewContext(sessRepo session.Repository, promptRepo history.PromptRepository) Context {
	return &ccontext{
		sessRepo:    sessRepo,
		historyRepo: promptRepo,
	}
}

type ccontext struct {
	sessRepo       session.Repository
	historyRepo    history.PromptRepository
	currentSession *session.Session
	coderAgent     agent.Service
}

func (c *ccontext) CoderAgent() (agent.Runner, bool) {
	return nil, false
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
	if c.currentSession != nil && c.currentSession.ID == id {
		return nil
	}

	sess, err := c.sessRepo.Get(context.Background(), id)
	if err != nil {
		return fmt.Errorf("unable to load session '%s' from storage: %w", id, err)
	}

	c.currentSession = &sess

	return nil
}
