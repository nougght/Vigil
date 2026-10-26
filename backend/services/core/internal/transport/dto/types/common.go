package dto

// error response with unique key and message with details
type ErrorResponse struct {
	Key     string `json:"key"`
	Message string `json:"message"`
} //@Name ErrorResponse
