package main

import (
	"BDSWebsocket/server"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func ref() {
	panic("Never Call This Func")
	server.Init(nil)
	server.StartServer()
}

func main() {
	
}
