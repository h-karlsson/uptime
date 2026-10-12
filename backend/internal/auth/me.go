package auth

import (
	"errors"
	"net/http"

	"github.com/h-karlsson/uptime/internal/user"
)

// Me returns the logged-in user. It must run behind RequireUser.
func Me(request *http.Request) (user.User, error) {
	loggedIn, ok := UserFromContext(request.Context())
	if !ok {
		return user.User{}, errors.New("me: no user in context; route is missing RequireUser")
	}
	return loggedIn, nil
}
