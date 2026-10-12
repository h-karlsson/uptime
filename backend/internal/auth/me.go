package auth

import (
	"net/http"

	"github.com/h-karlsson/uptime/internal/user"
)

func Me(request *http.Request) (user.User, error) {
	return CurrentUser(request.Context())
}
