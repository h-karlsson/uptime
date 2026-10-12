package monitor

import (
	"net/http"

	customHttp "github.com/h-karlsson/uptime/internal/http"
)

func monitorID(request *http.Request) (int64, error) {
	id, ok := customHttp.PathID(request, "id")
	if !ok {
		return 0, ErrNotFound
	}
	return id, nil
}
