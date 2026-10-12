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

// MinCost keeps the tests fast; the real service hashes with DefaultCost.
var demoPasswordHash = func() string {
	hash, err := bcrypt.GenerateFromPassword([]byte("demo1234"), bcrypt.MinCost)
	if err != nil {
		panic(err)
	}
	return string(hash)
}()

// fakeUserStore stands in for Postgres; demo@uptime.io / demo1234 is already registered.
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

type fakeSessionStore struct {
	users   *fakeUserStore
	userIDs map[string]int64
}

func (store *fakeSessionStore) Create(_ context.Context, tokenHash []byte, userID int64, _ time.Time) error {
	store.userIDs[string(tokenHash)] = userID
	return nil
}

func (store *fakeSessionStore) GetUser(_ context.Context, tokenHash []byte) (user.User, error) {
	userID, exists := store.userIDs[string(tokenHash)]
	if !exists {
		return user.User{}, session.ErrInvalidSession
	}
	for _, existing := range store.users.users {
		if existing.ID == userID {
			return existing, nil
		}
	}
	return user.User{}, session.ErrInvalidSession
}

func (store *fakeSessionStore) Delete(_ context.Context, tokenHash []byte) error {
	delete(store.userIDs, string(tokenHash))
	return nil
}

func newTestRouter() *http.ServeMux {
	users := &fakeUserStore{
		users:  map[string]user.User{"demo@uptime.io": {ID: 1, Email: "demo@uptime.io"}},
		hashes: map[string]string{"demo@uptime.io": demoPasswordHash},
	}
	sessions := &fakeSessionStore{users: users, userIDs: map[string]int64{}}

	router := http.NewServeMux()
	RegisterRoutes(router, user.NewService(users), session.NewService(sessions))
	return router
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
