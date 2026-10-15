package auth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"github.com/h-karlsson/uptime/internal/session"
	"github.com/h-karlsson/uptime/internal/user"
	"golang.org/x/crypto/bcrypt"
)

var demoPasswordHash = func() string {
	hash, err := bcrypt.GenerateFromPassword([]byte("demo1234"), bcrypt.MinCost)
	if err != nil {
		panic(err)
	}
	return string(hash)
}()

type fakeUserStore struct {
	users  map[string]user.User
	hashes map[string]string
}

func (store *fakeUserStore) Create(_ context.Context, email, passwordHash string) (user.User, error) {
	if _, exists := store.users[email]; exists {
		return user.User{}, user.ErrEmailTaken
	}
	newUser := user.User{ID: int64(len(store.users) + 1), Email: email}
	store.users[email] = newUser
	store.hashes[email] = passwordHash
	return newUser, nil
}

func (store *fakeUserStore) GetByEmail(_ context.Context, email string) (user.User, string, error) {
	existing, exists := store.users[email]
	if !exists {
		return user.User{}, "", user.ErrNotFound
	}
	return existing, store.hashes[email], nil
}

type fakeSession struct {
	userID    int64
	expiresAt time.Time
}

type fakeSessionStore struct {
	users    *fakeUserStore
	sessions map[string]fakeSession
}

func (store *fakeSessionStore) Create(_ context.Context, tokenHash []byte, userID int64, expiresAt time.Time) error {
	store.sessions[string(tokenHash)] = fakeSession{userID: userID, expiresAt: expiresAt}
	return nil
}

func (store *fakeSessionStore) Get(_ context.Context, tokenHash []byte) (session.Session, error) {
	stored, exists := store.sessions[string(tokenHash)]
	if !exists || !stored.expiresAt.After(time.Now()) {
		return session.Session{}, session.ErrInvalidSession
	}
	for _, existing := range store.users.users {
		if existing.ID == stored.userID {
			return session.Session{User: existing, ExpiresAt: stored.expiresAt}, nil
		}
	}
	return session.Session{}, session.ErrInvalidSession
}

func (store *fakeSessionStore) Extend(_ context.Context, tokenHash []byte, expiresAt time.Time) error {
	stored := store.sessions[string(tokenHash)]
	stored.expiresAt = expiresAt
	store.sessions[string(tokenHash)] = stored
	return nil
}

func (store *fakeSessionStore) Delete(_ context.Context, tokenHash []byte) error {
	delete(store.sessions, string(tokenHash))
	return nil
}

func (store *fakeSessionStore) DeleteForUser(_ context.Context, userID int64) error {
	for tokenHash, stored := range store.sessions {
		if stored.userID == userID {
			delete(store.sessions, tokenHash)
		}
	}
	return nil
}

func (store *fakeSessionStore) DeleteExpired(context.Context) (int64, error) {
	return 0, nil
}

func (store *fakeSessionStore) ageSessions(left time.Duration) {
	for tokenHash, stored := range store.sessions {
		stored.expiresAt = time.Now().Add(left)
		store.sessions[tokenHash] = stored
	}
}

func newTestRouter() *http.ServeMux {
	router, _ := newTestRouterWithSessions()
	return router
}

func newTestRouterWithSessions() (*http.ServeMux, *fakeSessionStore) {
	users := &fakeUserStore{
		users:  map[string]user.User{"demo@uptime.io": {ID: 1, Email: "demo@uptime.io"}},
		hashes: map[string]string{"demo@uptime.io": demoPasswordHash},
	}
	sessions := &fakeSessionStore{users: users, sessions: map[string]fakeSession{}}

	router := http.NewServeMux()
	RegisterRoutes(router, user.NewService(users), session.NewService(sessions))
	return router, sessions
}

func send(router http.Handler, method, path, body string, cookies ...*http.Cookie) *httptest.ResponseRecorder {
	request := httptest.NewRequest(method, path, strings.NewReader(body))
	for _, cookie := range cookies {
		request.AddCookie(cookie)
	}
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	return response
}

func sessionCookie(response *httptest.ResponseRecorder) *http.Cookie {
	for _, cookie := range response.Result().Cookies() {
		if cookie.Name == session.CookieName {
			return cookie
		}
	}
	return nil
}
