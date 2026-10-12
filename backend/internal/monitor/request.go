package monitor

import (
	"net/http"
	"strconv"

	"github.com/h-karlsson/uptime/internal/auth"
)

func currentUserID(request *http.Request) (int64, error) {
	loggedIn, err := auth.CurrentUser(request.Context())
	return loggedIn.ID, err
}

func monitorID(request *http.Request) (int64, error) {
	id, err := strconv.ParseInt(request.PathValue("id"), 10, 64)
	if err != nil {
		return 0, ErrNotFound
	}
	return id, nil
}
