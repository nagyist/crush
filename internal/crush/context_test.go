package crush_test

import (
	"context"
	"testing"

	"github.com/charmbracelet/crush/internal/crush"
	"github.com/charmbracelet/crush/internal/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockSessionRepo struct {
	createdSession session.Session
	getInvoked     bool
	toGetSession   session.Session
}

func (m *mockSessionRepo) Create(ctx context.Context, title string) (session.Session, error) {
	return m.createdSession, nil
}

func (m *mockSessionRepo) CreateTitleSession(ctx context.Context, parentSessionID string) (session.Session, error) {
	return session.Session{}, nil
}

func (m *mockSessionRepo) CreateTaskSession(ctx context.Context, toolCallID, parentSessionID, title string) (session.Session, error) {
	return session.Session{}, nil
}

func (m *mockSessionRepo) Get(ctx context.Context, id string) (session.Session, error) {
	m.getInvoked = true
	return m.toGetSession, nil
}

func (m *mockSessionRepo) List(ctx context.Context) ([]session.Session, error) {
	return []session.Session{}, nil
}

func (m *mockSessionRepo) Save(ctx context.Context, sess session.Session) (session.Session, error) {
	return session.Session{}, nil
}

func (m *mockSessionRepo) Delete(ctx context.Context, id string) error {
	return nil
}

func TestContext_ResolveCurrentSessionCreatesNewSession(t *testing.T) {
	justCreatedSession := session.Session{
		ID: "test-just-created-session",
	}
	sessRepo := mockSessionRepo{
		createdSession: justCreatedSession,
	}

	cctx := crush.NewContext(&sessRepo)

	sess, err := cctx.ResolveCurrentSession()
	require.NoError(t, err)
	assert.Equal(t, justCreatedSession, sess)
}

func TestContext_MakeSessionCurrentLoadsAndStoresDataForGivenSessionByID(t *testing.T) {
	sessionToGet := session.Session{
		ID: "test-old-dusty-session",
	}
	sessRepo := mockSessionRepo{
		toGetSession: sessionToGet,
	}

	cctx := crush.NewContext(&sessRepo)

	err := cctx.MakeSessionCurrent("test-old-dusty-session")
	require.NoError(t, err)
	assert.True(t, sessRepo.getInvoked)

	sess, err := cctx.ResolveCurrentSession()
	require.NoError(t, err)
	assert.Equal(t, sessionToGet, sess)
}
