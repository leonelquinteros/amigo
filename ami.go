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

// SendTo sends a PRIVMSG command to a user or a channel specified on 'dest' param.
func (a *Amigo) SendTo(dest, msg string) error {
    command := "PRIVMSG " + dest + " :" + msg

    return a.Send(command)
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
    /* Extra debug
    if msg.Prefix != nil {
        log.Println("Name: " + msg.Prefix.Name)
        log.Println("User: " + msg.Prefix.User)
        log.Println("Host: " + msg.Prefix.Host)
    }
    log.Println("Command: " + msg.Command)
    log.Println("Params: " + strings.Join(msg.Params, " "))
    log.Println("Trailing: " + msg.Trailing)
    */

    // Handle PING
    if msg.Command == "PING" {
        a.Send("PONG :" + msg.Trailing)
    }

    // Handle message
    if msg.Command == "PRIVMSG" {
        // Are you talking to me?
        if strings.HasPrefix(msg.Trailing, a.mem.Nick) {
            a.handleCommand(msg)
        } else {
            // Free talk
            a.handleConversation(msg)
        }
    }
}

// handleCommand parses and dispatches commands sent directly to the bot using the nick.
func (a *Amigo) handleCommand(msg *irc.Message) {
    // Parse command
    cmd, err := a.ParseCommand(msg)
    if err != nil {
        log.Println("AMIGO ERROR: " + err.Error())
        return
    }

    switch {
    case cmd.Method == "say":
        a.Say(cmd)

    case cmd.Method == "set master":
        a.SetMaster(cmd)
    }
}

func (a *Amigo) handleConversation(msg *irc.Message) {
    // Nothing yet...
}

// Say works like an Echo. Takes a Command and returns the params to the sender.
func (a *Amigo) Say(c *Command) {
    text := strings.TrimSpace(strings.Join(c.Params, " "))

    if text != "" {
        a.SendTo(c.Dest, text)
    }
}

func (a *Amigo) SetMaster(c *Command) {
    master := c.Params[0]

    if master == "" {
        return
    }

    if a.mem.Masters != nil {
        for _, m := range a.mem.Masters {
            if m == master {
                return
            }
        }
    }

    a.mem.Masters = append(a.mem.Masters, master)

    a.SendTo(master, "Welcome, master " + master)
}
