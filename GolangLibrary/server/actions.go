package server

import (
	"BDSWebsocket/protocol"
	"errors"
)

var ErrNoHandler = errors.New("no handler for action")

// HandlerMap is a map of action name to action handler
type HandlerMap map[string]func(client *Client, packet *protocol.PacketBase) (*protocol.PacketBase, error)

var RegisteredHandlers = make(HandlerMap)

func (h HandlerMap) Call(client *Client, packet *protocol.PacketBase) (*protocol.PacketBase, error) {
	handler, ok := h[packet.Action]
	if !ok {
		return nil, ErrNoHandler
	}
	return handler(client, packet)
}

// Register set the handler for Actions
func (h HandlerMap) Register() {
	h[protocol.LoginRequest_Action_Key] = HandleLogin
	h[protocol.RuncmdRequest_Action_Key] = HandleRuncmd
	h[protocol.BroadcastRequest_Action_Key] = HandleBroadcast
}
