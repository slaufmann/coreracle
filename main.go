// Copyright (C) 2020 Stefan Laufmann
//
// This file is part of coreracle.
//
// coreracle is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// coreracle is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with coreracle.  If not, see <https://www.gnu.org/licenses/>.

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
	client.HandleFunc(irc.DISCONNECTED,
						func (conn *irc.Conn, line *irc.Line) { quitSig <- true })
	client.HandleFunc(irc.CONNECTED, joinOnConnect)

	// connect!
	if err := client.Connect(); err != nil {
		fmt.Printf("Connection error: %s\n", err.Error())
	}

	// wait for disconnect
	<-quitSig
}

func joinOnConnect(conn *irc.Conn, line *irc.Line) {
	fmt.Printf(conn.String())
	channelArg := "#botwar"
	conn.Join(channelArg)
}
