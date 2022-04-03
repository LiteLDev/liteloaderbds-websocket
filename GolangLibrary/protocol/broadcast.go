package protocol

const (
	BroadcastRequest_Action_Key  = "BroadcastRequest"
	BroadcastResponse_Action_Key = "BroadcastResponse"
)

type BroadcastRequest struct {
	// Message that send to everyone
	Message string `json:"Message" validate:"required"`

	// MessageType of the message
	// Ref => Types.hpp->TextType
	// default: 1
	//    TextType::RAW = 0,
	//    TextType::CHAT = 1,
	//    TextType::TRANSLATION = 2,
	//    TextType::POPUP = 3,
	//    TextType::JUKEBOX_POPUP = 4,
	//    TextType::TIP = 5,
	//    TextType::SYSTEM = 6,
	//    TextType::WHISPER = 7,
	//    TextType::ANNOUNCEMENT = 8,
	//    TextType::JSON_WHISPER = 9,
	//    TextType::JSON = 10
	MessageType int `json:"MessageType" validate:"gte=0,lte=128"`
}

// BroadcastResponse is the response to BroadcastRequest that contains nothing
type BroadcastResponse struct {
}
