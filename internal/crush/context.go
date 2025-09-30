package crush

import (
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

func NewContext() Context {
	return &context{}
}

type context struct {
	sessions session.Service
}

func (c *context) ResolveCurrentSession() (session.Session, error) {
	/*
		if p.session.ID == "" {
			newSession, err := p.app.Sessions.Create(context.Background(), "New Session")
			if err != nil {
				return util.ReportError(err)
			}
			session = newSession
			cmds = append(cmds, util.CmdHandler(chat.SessionSelectedMsg(session)))
		}
	*/
	return session.Session{}, nil
}

func (c *context) MakeSessionCurrent(id string) error {
	return nil
}

func (c *context) CoderAgent() (agent.Service, bool) {
	return nil, false
}
