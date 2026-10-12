package monitor

import (
	"errors"
	"fmt"
	"net/http"

	customHttp "github.com/h-karlsson/uptime/internal/http"
	"github.com/h-karlsson/uptime/internal/validation"
)

func toClientError(err error) error {
	switch {
	case errors.Is(err, ErrNotFound):
		return &customHttp.ClientError{Message: err.Error(), Status: http.StatusNotFound}
	case errors.Is(err, ErrDuplicateURL):
		return &customHttp.ClientError{Message: err.Error(), Status: http.StatusConflict}
	case errors.Is(err, ErrLimitReached):
		message := fmt.Sprintf("you can have at most %d monitors", MaxPerUser)
		return &customHttp.ClientError{Message: message, Status: http.StatusConflict}
	}
	if validationError, ok := errors.AsType[*validation.Error](err); ok {
		return &customHttp.ClientError{Message: validationError.Message}
	}
	return err
}
