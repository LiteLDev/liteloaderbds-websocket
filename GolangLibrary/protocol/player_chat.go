package protocol

const (
	PlayerChatEvent_Action_Key = "PlayerChatEvent"
)

// PlayerChatEvent [Server -> Client]
type PlayerChatEvent struct {
	// Player is the player's Xbox ID that sent the message.
	Player string `json:"Player" validate:"required"`
	// Message is the message that the player sent.
	Message string `json:"Message"`
}
