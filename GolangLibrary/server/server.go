package server

import (
	"BDSWebsocket/logger"
	"log"
	"net/http"
)

import "C"

//export StartServer
func StartServer() {
	Config.LoadConfig("plugins/llws.json")
	go ClientHub.run()
	http.HandleFunc(Config.Endpoint, func(w http.ResponseWriter, r *http.Request) {
		conn, err := websocketUpgrader.Upgrade(w, r, nil)
		if err != nil {
			logger.Println("Client Upgrade:", err)
			return
		}
		client := &Client{hub: ClientHub, conn: conn, send: make(chan []byte, 256)}
		client.hub.register <- client
		logger.Printf("New connection establish: %s\n", client.conn.RemoteAddr().String())
		go client.writeLoop()
		go client.readLoop()
	})
	logger.Printf("Websocket Server started at %s", Config.ListenAddr)
	go func() {
		if Config.UsingTLS {
			err := http.ListenAndServeTLS(Config.ListenAddr, Config.CertFile, Config.KeyFile, nil)
			if err != nil {
				log.Fatal("ListenAndServeTLS: ", err)
			}
			logger.Printf("ListenAndServeTLS: %s\n", Config.ListenAddr)
		} else {
			err := http.ListenAndServe(Config.ListenAddr, nil)
			if err != nil {
				log.Fatal("ListenAndServe: ", err)
			}
			logger.Printf("ListenAndServe: %s\n", Config.ListenAddr)
		}
	}()
	RegisteredHandlers.Register()
}
