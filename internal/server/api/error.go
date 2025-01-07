package api

import "fmt"

type ErrorText string

const (
	ErrorResponseInternalServerError ErrorText = "internal server error"

	ErrorResponseBoardNotFound ErrorText = "board not found"
	ErrorResponseNoAccess      ErrorText = "no access"
)

func ErrorResponseValidationError(err error) ErrorText {
	return ErrorText(fmt.Sprintf("validation error: %v", err))
}

type ErrorResponse struct {
	Error ErrorText `json:"error"`
}
