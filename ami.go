package main

import (
	"errors"
    "log"
    "strings"

	"github.com/sorcix/irc"
)

// This is the bot
type Amigo struct {
    // Memory
    mem *Memory

	// Connection handler
	conn *irc.Conn
}

// EhAmigo starts the bot.
func (a *Amigo) EhAmigo() {
    a.mem = LoadMemory()

	// Connect
	err := a.connect()
	if err != nil {
		log.Fatal(err.Error())
	}

	// Start
	go a.init()
	a.listen()
}

// Send sends a raw IRC message over the active network stream.
func (a *Amigo) Send(msg string) error {
    log.Println("-> " + msg)

	return a.conn.Encode(irc.ParseMessage(msg))
}

// connect starts the IRC connection and stores the handler in conn.
func (a *Amigo) connect() error {
    log.Println("Connecting to " + a.mem.Host)

	c, err := irc.Dial(a.mem.Host)

	if err != nil {
		errMsg := "AMIGO ERROR: " + err.Error()
		return errors.New(errMsg)
	}

	a.conn = c

	return nil
}

// init sends IRC setup commands.
func (a *Amigo) init() {
	a.Send("NICK " + a.mem.Nick)
	a.Send("USER " + a.mem.Nick + " 0 * :amigo")
	a.Send("JOIN " + a.mem.Channel)
}

// listen gets all the network stream and dispatches the messages.
func (a *Amigo) listen() {
	for {
		msg, err := a.conn.Decode()
		if err != nil {
			log.Fatal("AMIGO ERROR: " + err.Error())
		}

        go a.handleMessage(msg)
	}
}

// handleMessage gets messages received on the IRC network and parses them to recognize commands.
func (a *Amigo) handleMessage(msg *irc.Message) {
    log.Println(msg.String())
    if msg.Prefix != nil {
        log.Println("Nick: " + msg.Prefix.Name)
        log.Println("User: " + msg.Prefix.User)
        log.Println("Host: " + msg.Prefix.Host)
    }
    log.Println("Command: " + msg.Command)
    log.Println("Params: " + strings.Join(msg.Params, " "))
    log.Println("Trailing: " + msg.Trailing)

    switch {
        case msg.Command == "PING":
            a.Send("PONG :" + msg.Trailing)
    }
}
