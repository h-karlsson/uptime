package stats

import (
	"net/http"

	"github.com/h-karlsson/uptime/internal/auth"
	customHttp "github.com/h-karlsson/uptime/internal/http"
)

func requestIDs(request *http.Request) (userID, monitorID int64, err error) {
	userID, err = auth.CurrentUserID(request.Context())
	if err != nil {
		return 0, 0, err
	}
	monitorID, ok := customHttp.PathID(request, "id")
	if !ok {
		return 0, 0, toClientError(ErrNotFound)
	}
	return userID, monitorID, nil
}
