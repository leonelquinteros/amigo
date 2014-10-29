package main

import (
	"errors"
	"github.com/sorcix/irc"
)

// This is the bot
type Amigo struct {
	// Connection params
	host, channel, nick, master string

	// Connection handler
	conn *irc.Conn
}

// EhAmigo starts the bot.
func (a *Amigo) EhAmigo(host, channel, nick, master string) {
	// Config set
	a.host = host
	a.channel = channel
	a.nick = nick
	a.master = master

	// Connect
	err := a.connect()
	if err != nil {
		println("AMIGO ABORT! " + err.Error())
		return
	}

	// Start
	go a.init()
	a.listen()
}

// Send sends a raw IRC message over the active network stream.
func (a *Amigo) Send(msg string) error {
	return a.conn.Encode(irc.ParseMessage(msg))
}

// connect starts the IRC connection and stores the handler in conn.
func (a *Amigo) connect() error {
	c, err := irc.Dial(a.host)

	if err != nil {
		errMsg := "AMIGO ERROR: " + err.Error()
		println(errMsg)
		return errors.New(errMsg)
	}

	a.conn = c

	return nil
}

// init sends IRC setup commands.
func (a *Amigo) init() {
	a.Send("NICK " + a.nick)
	a.Send("USER " + a.nick + " 0 * :amigo")
	a.Send("JOIN " + a.channel)
}

// listen gets all the network stream and dispatches the messages.
func (a *Amigo) listen() {
	for {
		msg, err := a.conn.Decode()

		if err != nil {
			println("AMIGO ERROR: " + err.Error())
			a.conn.Close()
			break
		}

		println(msg.String())
	}
}
