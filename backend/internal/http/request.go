package http

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

// ClientError describes invalid input that can safely be explained to the client.
type ClientError struct {
	Message string
	// Status is the HTTP status to respond with; zero means 400 Bad Request.
	Status int
}

func (err *ClientError) Error() string {
	return err.Message
}

// ReadJSON requires exactly one JSON value and rejects unknown object fields.
func ReadJSON(request *http.Request, destination any) error {
	decoder := json.NewDecoder(request.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(destination); err != nil {
		if errors.Is(err, io.EOF) {
			return &ClientError{Message: "Request body is required"}
		}
		return &ClientError{Message: "Invalid JSON body: check syntax, field names, and field types"}
	}

	var extra any
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		return &ClientError{Message: "Request body must contain exactly one JSON value"}
	}
	return nil
}
