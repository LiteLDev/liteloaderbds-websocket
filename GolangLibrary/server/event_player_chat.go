package server

import (
	"BDSWebsocket/logger"
	"BDSWebsocket/protocol"
	"encoding/json"
	"reflect"
)

import "C"

//export ChatEventBroadcast
func ChatEventBroadcast(playerName string, message string) {
	chatEvent := protocol.PlayerChatEvent{
		Player:  playerName,
		Message: message,
	}
	packet := protocol.PacketBase{
		Action:   protocol.PlayerChatEvent_Action_Key,
		PacketId: -1,
		Params:   &chatEvent,
	}

	data, err := json.Marshal(packet)
	if err != nil {
		logger.Error.Printf("Failed to marshal EventMessage (Type=%s):%s", reflect.TypeOf(chatEvent).Name(), err.Error())
		return
	}

	ClientHub.broadcast <- data
}
