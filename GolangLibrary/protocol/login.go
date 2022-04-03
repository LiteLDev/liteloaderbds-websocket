package protocol

const (
	LoginRequest_Action_Key  = "LoginRequest"
	LoginResponse_Action_Key = "LoginResponse"
)

type LoginRequest struct {
	// Password is considered as the access key to the server.
	// It would only need to send once per session.
	Password string `json:"Password" validate:"required"`
}

type LoginResponse struct {
	// Success return true if the login is successful.
	Success bool `json:"Success" validate:"required"`
	// Message is the error message if the login is not successful.
	// It would be empty if the login is successful.
	Message string `json:"Message"`
}
