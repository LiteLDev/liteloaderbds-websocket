package protocol

const (
	ErrorResponse_Action_Key = "ErrorResponse"
)

// ErrorResponse is the packet send by the server to the client when an error occurs.
type ErrorResponse struct {
	// Message is the error message
	Message string `json:"Message"`
}
