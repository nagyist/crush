package crush_test

import (
	"context"
	"errors"
	"testing"

	"github.com/charmbracelet/crush/internal/crush"
	"github.com/charmbracelet/crush/internal/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockSessionRepo struct {
	createInvoked bool
	createdSession session.Session
	createErr      error
	getInvoked     bool
	toGetSession   session.Session
	getErr error
}

func (m *mockSessionRepo) Create(ctx context.Context, title string) (session.Session, error) {
	m.createInvoked = true
	return m.createdSession, m.createErr
}

func (m *mockSessionRepo) CreateTitleSession(ctx context.Context, parentSessionID string) (session.Session, error) {
	return session.Session{}, nil
}

func (m *mockSessionRepo) CreateTaskSession(ctx context.Context, toolCallID, parentSessionID, title string) (session.Session, error) {
	return session.Session{}, nil
}

func (m *mockSessionRepo) Get(ctx context.Context, id string) (session.Session, error) {
	m.getInvoked = true
	return m.toGetSession, m.getErr
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

func TestContext_ResolveCurrentSessionCreatesNewSessionOnlyOnce(t *testing.T) {
	justCreatedSession := session.Session{
		ID: "test-just-created-session",
	}
	sessRepo := mockSessionRepo{
		createdSession: justCreatedSession,
	}

	cctx := crush.NewContext(&sessRepo)

	sess, err := cctx.ResolveCurrentSession()
	require.NoError(t, err)
	assert.True(t, sessRepo.createInvoked)
	assert.Equal(t, justCreatedSession, sess)

	// set this back
	sessRepo.createInvoked = false

	sess, err = cctx.ResolveCurrentSession()
	require.NoError(t, err)
	assert.False(t, sessRepo.createInvoked)
	assert.Equal(t, justCreatedSession, sess)
}

func TestContext_ResolveCurrentSessionCreatesNewSessionErrorsOnCreate(t *testing.T) {
	sessRepo := mockSessionRepo{
		createErr: errors.New("test failed to create session entry"),
	}

	cctx := crush.NewContext(&sessRepo)

	_, err := cctx.ResolveCurrentSession()
	require.EqualError(t, err, "failed to create session: test failed to create session entry")
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

func TestContext_MakeSessionCurrentLoadsDoesNotInteractWithRepoSecondTime(t *testing.T) {
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

	// set this back
	sessRepo.getInvoked = false

	sess, err := cctx.ResolveCurrentSession()
	require.NoError(t, err)
	assert.False(t, sessRepo.createInvoked)
	assert.Equal(t, sessionToGet, sess)

	err = cctx.MakeSessionCurrent("test-old-dusty-session")
	require.NoError(t, err)
	assert.False(t, sessRepo.getInvoked)
}

func TestContext_MakeSessionCurrentLoadsAndStoresDataForGivenSessionByIDErrorsOnGet(t *testing.T) {
	sessRepo := mockSessionRepo{
		getErr: errors.New("test failed to get session entry"),
	}

	cctx := crush.NewContext(&sessRepo)

	err := cctx.MakeSessionCurrent("test-old-dusty-session")
	require.EqualError(t, err, "unable to load session 'test-old-dusty-session' from storage: test failed to get session entry")
}
