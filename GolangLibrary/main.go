package main

import (
	"BDSWebsocket/server"
)

func ref() {
	panic("Never Call This Func")
	server.Init()
	server.StartServer()
}

func main() {

}
