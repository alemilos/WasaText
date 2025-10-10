package api

import "fmt"

const (
	InternalServerError = "Internal server error"
	InvalidRequestBody  = "Invalid Request Body"
)

func ErrorMessage(msg ...string) string {
	text := ""
	if len(msg) > 0 && msg[0] != "" {
		text = msg[0]
	}
	return fmt.Sprintf(`{"message":"%s"}`, text)
}
