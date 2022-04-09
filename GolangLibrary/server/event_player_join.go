package server

import (
	"BDSWebsocket/protocol"
	"BDSWebsocket/server/logger"
	"encoding/json"
	"reflect"
)

import "C"

//export JoinEventBroadcast
func JoinEventBroadcast(playerName string, XUID string, UUID string, ipAddress string, pos []float32, dimensionId int) {
	joinEvent := protocol.PlayerJoinEvent{
		Player:      playerName,
		XUID:        XUID,
		UUID:        UUID,
		IpAddress:   ipAddress,
		Position:    protocol.Vec3FromSlice(pos),
		DimensionId: dimensionId,
	}
	packet := protocol.PacketBase{
		Action:   protocol.PlayerJoinEvent_Action_Key,
		PacketId: -1,
		Params:   &joinEvent,
	}

	data, err := json.Marshal(packet)
	if err != nil {
		logger.Error.Printf("Failed to marshal EventMessage (Type=%s):%s", reflect.TypeOf(joinEvent).Name(), err.Error())
		return
	}

	ClientHub.broadcast <- data
}
