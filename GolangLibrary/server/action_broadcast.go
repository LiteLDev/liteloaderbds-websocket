package server

import (
	"BDSWebsocket/protocol"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"reflect"
)

func HandleBroadcast(client *Client, packet *protocol.PacketBase) (*protocol.PacketBase, error) {
	if !client.isAuthorized() {
		return nil, fmt.Errorf("client %s is not authorized", client.conn.RemoteAddr().String())
	}

	broadcastRequest := protocol.BroadcastRequest{
		MessageType: 1,
	}
	err := mapstructure.Decode(packet.Params, &broadcastRequest)
	if err != nil {
		return nil, fmt.Errorf("cannot cast to %s %w", reflect.TypeOf(broadcastRequest).Name(), err)
	}
	if ok := ValidateStructPrint(broadcastRequest); !ok {
		return nil, fmt.Errorf("cannot validate %s", reflect.TypeOf(broadcastRequest).Name())
	}
	CallBroadcastMessageWrapper(broadcastRequest.Message, broadcastRequest.MessageType)
	return &protocol.PacketBase{
		Action:   protocol.BroadcastResponse_Action_Key,
		PacketId: packet.PacketId,
		Params:   &protocol.BroadcastResponse{},
	}, nil
}
