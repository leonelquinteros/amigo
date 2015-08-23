package main

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
const param_delimiter = ";;"

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

type Command struct {
	Method  string
	Keyword string
	Params  []string
	Dest    string
}

// ParseCommand receives a command in a raw string and returns a new Command struct
func (a *Amigo) ParseCommand(msg *irc.Message) (*Command, error) {
	c := new(Command)

	// Remove nick from message
	raw := strings.TrimSpace(msg.Trailing[len(a.Nick):])

	// Command
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
	c.Params = strings.Split(raw, " "+param_delimiter)
	for key, param := range c.Params {
		// Preserve escaped param delimiters in text
		c.Params[key] = strings.TrimSpace(strings.Replace(param, "\\"+param_delimiter, param_delimiter, -1))
	}

	// Talking on a channel or private?
	receiver := strings.Join(msg.Params, " ")

	if receiver == a.Nick && msg.Prefix != nil {
		c.Dest = msg.Prefix.Name
	} else {
		// Talking on a channel
		c.Dest = receiver
	}

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
