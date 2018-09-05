package amigo

import (
	"encoding/json"
	"errors"
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/sorcix/irc"
)

// Amigo is the bot object
type Amigo struct {
	// Config
	Host, Channel, Nick, Password string

	// Memory
	mem *Memory

	// Connection handler
	conn *irc.Conn

	// Quit receiver
	quit chan bool
}

// EhAmigo starts the bot.
func (a *Amigo) EhAmigo(host, channel, nick, password string) {
	a.Host = host
	a.Channel = channel
	a.Nick = nick
	a.Password = password

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
func (a *Amigo) SendTo(dest, msg string) (err error) {
	clean := strings.Replace(msg, "\r\n", "\n", -1)
	lines := strings.Split(clean, "\n")

	for _, l := range lines {
		if l == "" {
			continue
		}

		command := "PRIVMSG " + dest + " :" + l
		err = a.Send(command)
		if err != nil {
			return err
		}

		time.Sleep(500 * time.Millisecond) // Prevent Excess Flood
	}

	return nil
}

// connect starts the IRC connection and stores the handler in conn.
func (a *Amigo) connect() error {
	log.Println("Connecting to " + a.Host)

	c, err := irc.Dial(a.Host)

	if err != nil {
		errMsg := "AMIGO ERROR: " + err.Error()
		return errors.New(errMsg)
	}

	a.conn = c

	return nil
}

// init sends IRC setup commands.
func (a *Amigo) init() {
	a.Send("NICK " + a.Nick)
	a.Send("USER " + a.Nick + " 0 * :amigo")
	a.Send("JOIN " + a.Channel)
}

// listen gets all the network stream and dispatches the messages.
func (a *Amigo) listen() {
Listen:
	for {
		select {
		case <-a.quit:
			a.Send("QUIT Shutting down")
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
	// Handle panics
	defer func() {
		if r := recover(); r != nil {
			log.Println("!!! I'm panicking !!! ", r)
		}
	}()

	// Log message
	log.Println(msg.String())

	// Handle PING
	if msg.Command == "PING" {
		a.Send("PONG :" + msg.Trailing)
	}

	// Handle message
	if msg.Command == "PRIVMSG" {
		// Are you talking to me?
		if strings.HasPrefix(msg.Trailing, a.Nick) {
			a.handleCommand(msg)
		} else {
			// Free talk ('say when', 'cmd when')
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

	a.dispatchCommand(cmd)
}

// dispatchCommand matches a command name to a method and executes it.
func (a *Amigo) dispatchCommand(cmd *Command) {
	switch {
	case cmd.Method == "help":
		a.Help(cmd)

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

	case cmd.Method == "exec when":
		var exec string

		if len(cmd.Params) > 1 {
			exec = cmd.Params[1]
		}

		a.ExecWhen(cmd.Params[0], exec)

	case cmd.Method == "cmd":
		var c string

		if len(cmd.Params) > 1 {
			c = strings.Join(cmd.Params[1:], " "+paramDelimiter)
		}

		a.DefineCommand(cmd.Params[0], c)

	case cmd.Method == "sys run":
		a.SysRun(cmd)

	}
}

func (a *Amigo) handleConversation(msg *irc.Message) {
	// Who to answer?
	who := a.getDestinatary(msg)

	// Check commands
	for k, v := range a.mem.AutoCmd {
		if strings.Contains(msg.Trailing, k) {
			// Create command
			c, err := a.getCommand(v)
			if err != nil {
				log.Println("AMIGO ERROR: " + err.Error())
				return
			}
			c.Dest = who

			a.dispatchCommand(c)
		}
	}
}

// Help displays help information about commands.
func (a *Amigo) Help(c *Command) {
	help := `-- Amigobot help --
I can be controlled through commands sent to me. 
I'll will only answer to commands that starts with my nick in the form: 
[NICK] [COMMAND] [[PARAM_DELIMITER] [PARAM]]...
Commands:
- help: This very command.
- say: Makes me say something.
- tell me: Display information about myself.
- set master: Sets a new master.
- del master: Deletes an existent master.
- set password: Changes the master password.
- join: Makes me join a channel.
- leave: Makes me leave a channel.
- shutdown: Makes me to kill myself.
- cmd: Defines a new command.
- exec when: Defines a command to execute when somebody says something.
- sys run: Runs a command in the local system (DANGER)
`

	a.SendTo(c.Dest, help)
}

// Say works like an Echo. Takes a Command and returns the params to the sender.
func (a *Amigo) Say(c *Command) {
	text := strings.TrimSpace(strings.Join(c.Params, " "))

	if text != "" {
		a.SendTo(c.Dest, text)
	}
}

// Tell is meant to display memory content to the IRC conversation.
// Masters:
//      tell masters
//      tell your masters
//
//  - will display master nicks loaded to memory.
func (a *Amigo) Tell(c *Command) {
	what := strings.ToLower(strings.TrimSpace(strings.Join(c.Params, " ")))

	switch {
	case what == "masters" || what == "your masters":
		a.SendTo(c.Dest, strings.Join(a.mem.Masters, ","))

	case what == "memory":
		raw, _ := json.MarshalIndent(a.mem, "", "  ")
		a.SendTo(c.Dest, string(raw))

	default:
		// Echo
		a.SendTo(c.Dest, what)
	}
}

// SetMaster adds a new master nick to memory.
// A master nick is a user who doesn't needs to pass the password parameter to exec a command to the bot.
// Is recommended that the first action to the bot is to set the first master to avoid spreading the password on public IRC conversations.
func (a *Amigo) SetMaster(c *Command) {
	if len(c.Params) < 1 {
		return
	}
	if c.Params[0] == "" {
		return
	}

	if a.mem.Masters != nil {
		for _, m := range a.mem.Masters {
			if m == c.Params[0] {
				return
			}
		}
	}

	a.mem.Masters = append(a.mem.Masters, c.Params[0])

	a.SendTo(c.Params[0], "Welcome, master "+c.Params[0])
}

// DelMaster will remove a master nick from memory.
func (a *Amigo) DelMaster(c *Command) {
	if len(c.Params) < 1 {
		return
	}
	if c.Params[0] == "" {
		return
	}

	if a.mem.Masters != nil {
		for key, m := range a.mem.Masters {
			if m == c.Params[0] {
				a.mem.Masters = append(a.mem.Masters[:key], a.mem.Masters[key+1:]...)
				break
			}
		}
	}

	a.SendTo(c.Params[0], "Goodbye "+c.Params[0]+", i'm not listening to you anymore")
}

// SetNick will make the bot to change its nick.
// Be careful with Nick Servers and already reserved nick names.
func (a *Amigo) SetNick(c *Command) {
	if len(c.Params) < 1 {
		return
	}
	if c.Params[0] == "" {
		return
	}

	a.Nick = c.Params[0]
	a.Send("NICK " + c.Params[0])
}

// SetPassword will change the master password on the bots memory.
func (a *Amigo) SetPassword(c *Command) {
	if len(c.Params) < 1 {
		return
	}
	if c.Params[0] == "" {
		return
	}

	a.Password = c.Params[0]

	a.SendTo(c.Dest, "Password changed")
}

// Join the IRC Channel provided.
func (a *Amigo) Join(c *Command) {
	if len(c.Params) < 1 {
		return
	}
	if c.Params[0] == "" {
		return
	}

	a.Send("JOIN " + c.Params[0])
}

// Leave the IRC channel provided.
func (a *Amigo) Leave(c *Command) {
	if len(c.Params) < 1 {
		return
	}
	if c.Params[0] == "" {
		return
	}

	a.Send("PART " + c.Params[0])
}

// Shutdown will gracefully disconnect from the IRC server and terminate the running process on the host machine.
func (a *Amigo) Shutdown() {
	a.quit <- true
}

// DefineCommand takes a new keyword for an existent command and uses it as an alias.
func (a *Amigo) DefineCommand(keyword, command string) {
	// Delete if empty
	if command == "" {
		if _, ok := a.mem.Commands[keyword]; ok {
			delete(a.mem.Commands, keyword)
		}

		return
	}

	// Set new
	a.mem.Commands[keyword] = command
}

// ExecWhen defines commands to execute when somebody in the channel (or private message) says something.
func (a *Amigo) ExecWhen(when, exec string) {
	// Delete if empty
	if exec == "" {
		if _, ok := a.mem.AutoCmd[when]; ok {
			delete(a.mem.AutoCmd, when)
		}

		return
	}

	// Set new
	a.mem.AutoCmd[when] = exec
}

// SysRun will execute the command provided on the host machine.
// When possible, it will return the command output to the IRC channel.
func (a *Amigo) SysRun(c *Command) {
	if len(c.Params) < 1 {
		return
	}

	if c.Params[0] == "" {
		return
	}

	var sysCmd *exec.Cmd

	if len(c.Params) > 1 {
		sysCmd = exec.Command(c.Params[0], c.Params[1:]...)
	} else {
		sysCmd = exec.Command(c.Params[0])
	}

	out, err := sysCmd.Output()
	if err != nil {
		a.SendTo(c.Dest, "SYS ERROR: "+err.Error())
	}

	a.SendTo(c.Dest, string(out))
}
