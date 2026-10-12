package validation

// Error describes invalid input that can safely be explained to the client.
type Error struct {
	Message string
}

func (err *Error) Error() string {
	return err.Message
}
