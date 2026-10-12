package http

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
)

type JSONHandler func(request *http.Request) (any, error)

const applicationJson = "application/json"
const contentType = "Content-Type"

func FormatResponseToJson(responseWriter http.ResponseWriter, status int, data any) {
	responseWriter.Header().Set(contentType, applicationJson)
	responseWriter.WriteHeader(status)

	if err := json.NewEncoder(responseWriter).Encode(data); err != nil {
		slog.Error("write JSON response", "error", err)
	}
}

func HandleJSON[T any](handler func(*http.Request) (T, error)) http.HandlerFunc {
	return HandleJSONStatus(http.StatusOK, handler)
}

func HandleJSONStatus[T any](status int, handler func(*http.Request) (T, error)) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		data, err := handler(request)
		if err != nil {
			WriteError(response, err)
			return
		}

		FormatResponseToJson(response, status, data)
	}
}

func HandleNoContent(handler func(*http.Request) error) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		if err := handler(request); err != nil {
			WriteError(response, err)
			return
		}

		response.WriteHeader(http.StatusNoContent)
	}
}

// WriteError responds with a ClientError's message and status, or a generic 500 for any other error.
func WriteError(response http.ResponseWriter, err error) {
	var clientError *ClientError
	if errors.As(err, &clientError) {
		status := clientError.Status
		if status == 0 {
			status = http.StatusBadRequest
		}
		FormatResponseToJson(response, status, map[string]string{"error": clientError.Message})
		return
	}
	slog.Error("request failed", "error", err)
	FormatResponseToJson(
		response,
		http.StatusInternalServerError,
		map[string]string{"error": "Internal server error"},
	)
}
