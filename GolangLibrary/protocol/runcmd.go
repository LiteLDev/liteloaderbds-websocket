package protocol

const (
	RuncmdRequest_Action_Key  = "RuncmdRequest"
	RuncmdResponse_Action_Key = "RuncmdResponse"
)

type RuncmdRequest struct {
	// Command to run
	Command string `json:"Command" validate:"required"`
}

type RuncmdResponse struct {
	// Message is the command message
	Message string `json:"Message"`
}
