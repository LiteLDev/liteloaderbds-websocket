package server

import (
	"BDSWebsocket/protocol"
	"BDSWebsocket/server/logger"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"reflect"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var websocketUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type Client struct {
	hub  *Hub
	conn *websocket.Conn
	auth bool
	send chan []byte
}

func (c *Client) isAuthorized() bool {
	return c.auth
}

func (c *Client) sendErrorFeedback(packetBase *protocol.PacketBase, msg string) {
	errPacket := &protocol.ErrorResponse{
		Message: msg,
	}
	outcomePacketBase := &protocol.PacketBase{
		Action:   protocol.ErrorResponse_Action_Key,
		PacketId: packetBase.PacketId,
		Params:   errPacket,
	}
	data, err := json.Marshal(outcomePacketBase)
	if err != nil {
		logger.Error.Printf("Failed to send error feedback (ID=%d,Client:%s):%s", packetBase.PacketId, c.conn.RemoteAddr().String(), err.Error())
		return
	}
	c.send <- data
}

func (c *Client) handlePacket(packet []byte) {
	var err error
	packetBase := protocol.PacketBase{}
	err = json.Unmarshal(packet, &packetBase)
	if err != nil {
		logger.Error.Println(err)
		return
	}
	if err := ValidateStructPrint(packetBase); !err {
		return
	}
	returnVal, err := RegisteredHandlers.Call(c, &packetBase)
	if err != nil {
		logger.Error.Println(err)
		c.sendErrorFeedback(&packetBase, err.Error())
		return
	}
	data, err := json.Marshal(returnVal)
	if err != nil {
		logger.Error.Printf("Failed to marshal response (ID=%d,Type=%s,Client:%s):%s", packetBase.PacketId, reflect.TypeOf(returnVal).Name(), c.conn.RemoteAddr().String(), err.Error())
		return
	}
	c.send <- data

}

func (c *Client) readLoop() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		logger.Printf("Received message [%s]: %s", c.conn.RemoteAddr(), string(message))
		c.handlePacket(message)
	}
}

func (c *Client) writeLoop() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			err := c.conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				return
			}

			time.Sleep(time.Second * 1)
			// Add rest of queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				err := c.conn.WriteMessage(websocket.TextMessage, <-c.send)
				if err != nil {
					return
				}
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
