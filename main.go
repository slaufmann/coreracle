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
	"os"
	"strings"

	irc "github.com/fluffle/goirc/client"
	"github.com/akamensky/argparse"
)

type options struct {	nick string
						channel string
						server string
						port string
					}

var opts = options{	nick: "coreracleBot",
					channel: "#botwar",
					server: "irc.freenode.net",
					port: "7000"}

func main() {
	// parse command line arguments
	parser := argparse.NewParser("coreracle", "Helpful IRC bot that can tell you the future based on coredumps and stacktraces.")
	var nickArg *string = parser.String("n", "nickname",
											&argparse.Options{Required: false, Help: "nick with which the bot joins a channel"})
	var chanArg *string = parser.String("c", "channel",
											&argparse.Options{Required: false, Help: "channel the bot should join"})
	var serverArg *string = parser.String("s", "server",
											&argparse.Options{Required: false, Help: "server the bot should connect to"})
	var portArg *string = parser.String("p", "port",
											&argparse.Options{Required: false, 
																Help: "port that should be used to connecto to the server"})
	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		os.Exit(1)
	}
	if (*nickArg != "") {
		opts.nick = *nickArg
	}
	if (*chanArg != "") {
		opts.channel = *chanArg
	}
	if (*serverArg != "") {
		opts.server = *serverArg
	}
	if (*portArg != "") {
		opts.port = *portArg
	}

	// create config and adjust settings
	config := irc.NewConfig(opts.nick)
	config.SSL = true
	config.SSLConfig = &tls.Config{ServerName: opts.server}
	config.Server = opts.server + ":" + opts.port
	config.NewNick = func(n string) string {return n + "^" }

	// create the client
	client := irc.Client(config)

	// disconnect signal
	quitSig := make(chan bool)

	// register handler functions
	client.HandleFunc(irc.DISCONNECTED,
						func (conn *irc.Conn, line *irc.Line) { quitSig <- true })
	client.HandleFunc(irc.CONNECTED, joinOnConnect)
	client.HandleFunc(irc.PRIVMSG, handlePrivMsg)

	// connect!
	if err := client.Connect(); err != nil {
		fmt.Printf("Connection error: %s\n", err.Error())
	}

	// wait for disconnect
	<-quitSig
}

func joinOnConnect(conn *irc.Conn, line *irc.Line) {
	fmt.Printf(conn.String())
	conn.Join(opts.channel)
}

func handlePrivMsg(conn *irc.Conn, line *irc.Line) {
	text := line.Text()
	strPart := strings.Split(text, ":")
	if (len(strPart[0]) != len(text)) {
		fmt.Printf("Someone was highlighted, maybe.\n")
		nickStruct := conn.Me()
		nick := nickStruct.Nick
		if (strPart[0] == nick) {
			fmt.Printf("It's us!\n")
		}
	}

	if line.Public() {
		fmt.Printf("Public message: %s\n",text)
		replyToMsg(conn, line.Target(), text)
	} else {
		fmt.Printf("Private message: %s\n", text)
		replyToMsg(conn, line.Target(), text)
	}
}

func replyToMsg(conn *irc.Conn, target string, text string) {
	fmt.Printf("Replying with message: %s\n", text)
	conn.Privmsg(target, text)
}
