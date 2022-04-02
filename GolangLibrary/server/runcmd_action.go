package server

import (
	"BDSWebsocket/protocol"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"reflect"
)

func HandleRuncmd(client *Client, packet *protocol.PacketBase) (*protocol.PacketBase, error) {
	if !client.isAuthorized() {
		return nil, fmt.Errorf("client %s is not authorized", client.conn.RemoteAddr().String())
	}

	loginRequest := protocol.RuncmdRequest{}
	err := mapstructure.Decode(packet.Params, &loginRequest)
	if err != nil {
		return nil, fmt.Errorf("cannot cast to %s %w", reflect.TypeOf(loginRequest).Name(), err)
	}
	if ok := ValidateStructPrint(loginRequest); !ok {
		return nil, fmt.Errorf("cannot validate %s", reflect.TypeOf(loginRequest).Name())
	}
	result := CallRuncmdFunc(loginRequest.Command)
	return &protocol.PacketBase{
		Action:   protocol.RuncmdResponse_Action_Key,
		PacketId: packet.PacketId,
		Params: &protocol.RuncmdResponse{
			Message: result,
		},
	}, nil
}
