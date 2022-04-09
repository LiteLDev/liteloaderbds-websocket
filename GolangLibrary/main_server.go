//go:build server

package main

import (
	"BDSWebsocket/server"
	"bufio"
	"fmt"
	"os"
	"strings"
)

func ref() {
	panic("Never Call This Func")
	server.Init()
	server.StartServer()
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Simple Shell")
	fmt.Println("---------------------")

	for {
		fmt.Print("-> ")
		text, _ := reader.ReadString('\n')
		// convert CRLF to LF
		text = strings.Replace(text, "\n", "", -1)

	}
}
