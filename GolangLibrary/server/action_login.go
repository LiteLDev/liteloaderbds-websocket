package server

import (
	"BDSWebsocket/protocol"
	"BDSWebsocket/server/logger"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"reflect"
)

func HandleLogin(client *Client, packet *protocol.PacketBase) (*protocol.PacketBase, error) {
	loginRequest := protocol.LoginRequest{}
	err := mapstructure.Decode(packet.Params, &loginRequest)
	if err != nil {
		return nil, fmt.Errorf("cannot cast to %s %w", reflect.TypeOf(loginRequest).Name(), err)
	}
	if ok := ValidateStructPrint(loginRequest); !ok {
		return nil, fmt.Errorf("cannot validate %s", reflect.TypeOf(loginRequest).Name())
	}
	if loginRequest.Password == Config.Token {
		client.auth = true
		logger.Printf("Client %s authenticated", client.conn.RemoteAddr().String())
	} else {
		client.auth = false
		logger.Printf("Client %s failed authentication", client.conn.RemoteAddr().String())
	}
	return &protocol.PacketBase{
		Action:   protocol.LoginResponse_Action_Key,
		PacketId: packet.PacketId,
		Params: &protocol.LoginResponse{
			Success: client.auth,
		},
	}, nil
}
