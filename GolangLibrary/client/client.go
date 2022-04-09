package client

import (
	"BDSWebsocket/protocol"
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"
	"log"
	"os"
	"strings"
)

var packetId int64 = 1

func sendLogin(conn *websocket.Conn, password string) {
	loginPacket := protocol.LoginRequest{Password: password}
	packetBase := protocol.PacketBase{PacketId: packetId, Action: protocol.LoginRequest_Action_Key, Params: loginPacket}
	data, _ := json.Marshal(packetBase)
	conn.WriteMessage(websocket.TextMessage, data)
}

func sendCommand(conn *websocket.Conn, command string) {
	commandPacket := protocol.RuncmdRequest{Command: command}
	packetBase := protocol.PacketBase{PacketId: packetId, Action: protocol.RuncmdRequest_Action_Key, Params: commandPacket}
	data, _ := json.Marshal(packetBase)
	conn.WriteMessage(websocket.TextMessage, data)
}

func ClientMain() {

	var serverAddr string
	var password string

	fmt.Println("Websocket Shell")

	fmt.Print("Server Address: ")
	fmt.Scanln(&serverAddr)

	fmt.Print("Password:       ")
	fmt.Scanln(&password)

	log.Printf("connecting to %s", serverAddr)

	c, _, err := websocket.DefaultDialer.Dial(serverAddr, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	go func() {
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				os.Exit(1)
			}
			packetBase := protocol.PacketBase{}
			json.Unmarshal(message, &packetBase)
			switch packetBase.Action {

			case protocol.LoginResponse_Action_Key:
				loginResponse := protocol.LoginResponse{}
				mapstructure.Decode(packetBase.Params, &loginResponse)
				if !loginResponse.Success {
					panic("Login failed")
				}

			case protocol.RuncmdResponse_Action_Key:
				runcmdResponse := protocol.RuncmdResponse{}
				mapstructure.Decode(packetBase.Params, &runcmdResponse)
				fmt.Printf("RuncmdResponse [PacketId=%d]\n%s\n", packetBase.PacketId, runcmdResponse.Message)

			case protocol.ErrorResponse_Action_Key:
				errorResponse := protocol.ErrorResponse{}
				mapstructure.Decode(packetBase.Params, &errorResponse)
				fmt.Printf("ErrorResponse [PacketId=%d]\n%s\n", packetBase.PacketId, errorResponse.Message)

			case protocol.PlayerJoinEvent_Action_Key:
				playerJoinEvent := protocol.PlayerJoinEvent{}
				mapstructure.Decode(packetBase.Params, &playerJoinEvent)
				fmt.Printf("PlayerJoinEvent [PacketId=%d]\n%#v\n", packetBase.PacketId, playerJoinEvent)
			case protocol.PlayerLeftEvent_Action_Key:
				playerLeftEvent := protocol.PlayerLeftEvent{}
				mapstructure.Decode(packetBase.Params, &playerLeftEvent)
				fmt.Printf("PlayerLeftEvent [PacketId=%d]\n%#v\n", packetBase.PacketId, playerLeftEvent)
			case protocol.PlayerChatEvent_Action_Key:
				playerChatEvent := protocol.PlayerChatEvent{}
				mapstructure.Decode(packetBase.Params, &playerChatEvent)
				fmt.Printf("PlayerChatEvent [PacketId=%d]\n%#v\n", packetBase.PacketId, playerChatEvent)

			}
		}
	}()

	sendLogin(c, password)

	reader := bufio.NewReader(os.Stdin)
	for {
		packetId++
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)
		text = strings.Replace(text, "\r", "", -1)
		sendCommand(c, text)
	}

}
