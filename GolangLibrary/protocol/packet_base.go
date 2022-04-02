package protocol

type PacketBase struct {
	// Action Specific the action type
	// This field is required and should not be empty
	Action string `json:"Action" validate:"required"`

	// PacketId Specific the packet id that client provide and server will return
	// it will be -1 if the packet is not a response
	// This field is required and should not be empty
	PacketId int64 `json:"PacketId" validate:"required"`

	// Params of the action
	Params interface{} `json:"Params"`
}
