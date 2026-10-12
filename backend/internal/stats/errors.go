package stats

import (
	"errors"
	"net/http"

	customHttp "github.com/h-karlsson/uptime/internal/http"
	"github.com/h-karlsson/uptime/internal/validation"
)

func toClientError(err error) error {
	if errors.Is(err, ErrNotFound) {
		return &customHttp.ClientError{Message: err.Error(), Status: http.StatusNotFound}
	}
	if validationError, ok := errors.AsType[*validation.Error](err); ok {
		return &customHttp.ClientError{Message: validationError.Message}
	}
	return err
}
