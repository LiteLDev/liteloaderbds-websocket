package server

import (
	"BDSWebsocket/protocol"
	"BDSWebsocket/server/logger"
	"encoding/json"
	"reflect"
)

import "C"

//export LeftEventBroadcast
func LeftEventBroadcast(playerName string, XUID string, UUID string, pos []float32, dimensionId int) {
	leftEvent := protocol.PlayerLeftEvent{
		Player:      playerName,
		XUID:        XUID,
		UUID:        UUID,
		Position:    protocol.Vec3FromSlice(pos),
		DimensionId: dimensionId,
	}
	packet := protocol.PacketBase{
		Action:   protocol.PlayerLeftEvent_Action_Key,
		PacketId: -1,
		Params:   &leftEvent,
	}

	data, err := json.Marshal(packet)
	if err != nil {
		logger.Error.Printf("Failed to marshal EventMessage (Type=%s):%s", reflect.TypeOf(leftEvent).Name(), err.Error())
		return
	}

	ClientHub.broadcast <- data
}
