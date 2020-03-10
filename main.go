package main

import (
	"crypto/tls"
	"fmt"
	irc "github.com/fluffle/goirc/client"
)

func main() {
	nickArg := "coreracleBot"
	serverArg := "irc.freenode.net"
	portArg := "7000"

	// create config and adjust settings
	config := irc.NewConfig(nickArg)
	config.SSL = true
	config.SSLConfig = &tls.Config{ServerName: serverArg}
	config.Server = serverArg + ":" + portArg
	config.NewNick = func(n string) string {return n + "^" }

	// create the client
	client := irc.Client(config)

	// disconnect signal
	quitSig := make(chan bool)

	// register handler functions
//	client.HandleFunc(irc.DISCONNECTED, quitOnDisconnect)

	// connect!
	if err := client.Connect(); err != nil {
		fmt.Printf("Connection error: %s\n", err.Error())
	}

	// wait for disconnect
	<-quitSig
}

//func quitOnDisconnect(conn *irc.Conn, line *irc.Line) {
//	quitSig <- true
//}

func joinOnConnect(conn *irc.Conn, line *irc.Line) {
	channelArg := "#afra"
	conn.Join(channelArg)
}
