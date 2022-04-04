package protocol

const (
	PlayerLeftEvent_Action_Key = "PlayerLeftEvent"
)

// PlayerLeftEvent [Server -> Client]
type PlayerLeftEvent struct {
	// Player is the player's Xbox ID that sent the message.
	Player string `json:"Player" validate:"required"`

	XUID string `json:"XUID" validate:"required"`

	UUID string `json:"UUID" validate:"required"`

	Position McVec3 `json:"Position" validate:"required"`

	DimensionId int `json:"DimensionId" validate:"required"`
}
