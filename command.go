package amigo

import (
	"errors"
	"github.com/sorcix/irc"
	"strings"
)

// Parameter delimiter string.
// Is used to separate different parameters for commands.
// The delimiter has to be preceded by a space (" ") character.
// If the delimiter string needs to be used as part of some data,
// can be escaped preceding a backslash (\) on front, like: "Some\;;data, not param delimiter here"
const paramDelimiter = ";;"

// Protocol commands list
var protocol = []string{
	"help",
	"tell me",
	"set master",
	"del master",
	"set nick",
	"set password",
	"join",
	"leave",
	"shutdown",
	"cmd",
	"say",
	"say when",
	"exec when",
	"sys run", // Fucking dangerous
}

// Command stores an Amigo command with all its different parts.
type Command struct {
	Method  string
	Keyword string
	Params  []string
	Dest    string
}

// ParseCommand receives a command in a raw string and returns a new Command struct
//
// Command format has to be in the following form to be recognized:
//
// [ NICK ] [ COMMAND ] [ PARAM ] [ SPACE + PARAM_DELIMITER [ PARAM ] ... ]
//
// Where NICK is the bot's nick on the IRC server,
// COMMAND is a string present in the protocol variable defined on this file,
// PARAM is any string formed by any chars other than a combination of SPACE and PARAM_DELIMITER,
// SPACE is the space character and PARAM_DELIMITER is defined in param_delimiter constant on command.go, a double semicolon ";;" by default.
//
// If PARAM_DELIMITER has to be included as part of a param value (usually not needed), you can escape it using a backslash \,
// Any combination of backslash and the param delimiter will be replaced by the param delimiter itself at parsing time.
//
func (a *Amigo) ParseCommand(msg *irc.Message) (*Command, error) {
	// Remove nick from message
	raw := strings.TrimSpace(msg.Trailing[len(a.Nick):])

	// Get command
	c, err := a.getCommand(raw)
	if err != nil {
		return nil, err
	}

	// Get response destinatary
	c.Dest = a.getDestinatary(msg)

	// Check auth
	if msg.Prefix == nil || msg.Prefix.Name == "" {
		return nil, errors.New("Empty user not authorized.")
	}
	auth := false
	if a.mem.Masters != nil {
		for _, master := range a.mem.Masters {
			if msg.Prefix.Name == master {
				auth = true
				break
			}
		}
	}
	if a.Password == c.Params[len(c.Params)-1] {
		auth = true
		c.Params = c.Params[:len(c.Params)-1]
	}
	if !auth {
		return nil, errors.New("User not authorized.")
	}

	return c, nil
}

// getCommand takes a string received on chat and tries to create a Command object from it.
func (a *Amigo) getCommand(raw string) (*Command, error) {
	// Command
	c := new(Command)

	found := false
	for _, cmd := range protocol {
		if strings.HasPrefix(raw, cmd) {
			found = true
			c.Method = strings.ToLower(cmd)
			c.Keyword = cmd
			break
		}
	}
	if !found {
		// Custom command
		for keyword, cmd := range a.mem.Commands {
			if strings.HasPrefix(raw, keyword) {
				found = true
				c.Method = strings.ToLower(cmd)
				c.Keyword = keyword
				break
			}
		}
	}
	if !found {
		return nil, errors.New("Command not found")
	}

	// Remove command from message.
	raw = strings.TrimSpace(raw[len(c.Keyword):])

	// Split params
	c.Params = strings.Split(raw, " "+paramDelimiter)
	for key, param := range c.Params {
		// Preserve escaped param delimiters in text
		c.Params[key] = strings.TrimSpace(strings.Replace(param, "\\"+paramDelimiter, paramDelimiter, -1))
	}

	return c, nil
}

// getDestinatary figures out if we're talking on a channel or in private.
// Returns the corresponding identifier to be used to answer.
func (a *Amigo) getDestinatary(msg *irc.Message) string {
	var d string
	receiver := strings.Join(msg.Params, " ")

	// Talking on a channel or private?
	if receiver == a.Nick && msg.Prefix != nil {
		d = msg.Prefix.Name
	} else {
		// Talking on a channel
		d = receiver
	}

	return d
}
