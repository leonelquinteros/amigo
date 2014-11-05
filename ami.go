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

    // Quit receiver
    quit chan bool
}

// EhAmigo starts the bot.
func (a *Amigo) EhAmigo() {
    a.quit = make(chan bool)

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
Listen:
	for {
        select {
            case <- a.quit:
                break Listen

            default:
                msg, err := a.conn.Decode()
                if err != nil {
                    log.Fatal("CONNECTION ERROR: " + err.Error())
                }

                go a.handleMessage(msg)
        }

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

    case cmd.Method == "tell me":
        a.Tell(cmd)

    case cmd.Method == "set master":
        a.SetMaster(cmd)

    case cmd.Method == "del master":
        a.DelMaster(cmd)

    case cmd.Method == "set nick":
        a.SetNick(cmd)

    case cmd.Method == "set password":
        a.SetPassword(cmd)

    case cmd.Method == "join":
        a.Join(cmd)

    case cmd.Method == "leave":
        a.Leave(cmd)

    case cmd.Method == "shutdown":
        a.Shutdown()
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

func (a *Amigo) Tell(c *Command) {
    what := strings.ToLower(strings.TrimSpace(strings.Join(c.Params, " ")))

    switch {
        case what == "masters" || what == "your masters":
            a.SendTo(c.Dest, strings.Join(a.mem.Masters, ","))
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

func (a *Amigo) DelMaster(c *Command) {
    master := c.Params[0]

    if master == "" {
        return
    }

    if a.mem.Masters != nil {
        for key, m := range a.mem.Masters {
            if m == master {
                a.mem.Masters = append(a.mem.Masters[:key], a.mem.Masters[key+1:]...)
                break
            }
        }
    }

    a.SendTo(master, "Goodbye " + master + ", i'm not listening to you anymore")
}

func (a *Amigo) SetNick(c *Command) {
    nick := c.Params[0]

    if nick == "" {
        return
    }

    a.mem.Nick = nick
    a.Send("NICK " + nick)
}


func (a *Amigo) SetPassword(c *Command) {
    password := c.Params[0]

    if password == "" {
        return
    }

    a.mem.Password = password

    a.SendTo(c.Dest, "Password changed")
}

func (a *Amigo) Join(c *Command) {
    channel := c.Params[0]

    if channel == "" {
        return
    }

    a.Send("JOIN " + channel)
}

func (a *Amigo) Leave(c *Command) {
    channel := c.Params[0]

    if channel == "" {
        return
    }

    a.Send("PART " + channel)
}

func (a *Amigo) Shutdown() {
    a.Send("QUIT Shutting down")

    a.quit <- true
}
