package protocol

const (
	PlayerJoinEvent_Action_Key = "PlayerJoinEvent"
)

// PlayerJoinEvent [Server -> Client]
type PlayerJoinEvent struct {
	// Player is the player's Xbox ID that sent the message.
	Player string `json:"Player" validate:"required"`

	XUID string `json:"XUID" validate:"required"`

	UUID string `json:"UUID" validate:"required"`

	Position McVec3 `json:"Position" validate:"required"`

	DimensionId int `json:"DimensionId" validate:"required"`

	IpAddress string `json:"IpAddress" validate:"required"`
}
