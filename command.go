package main

import (
	"errors"
	"github.com/sorcix/irc"
	"strings"
)

var param_delimiter = ";;"

// Protocol commands list
var protocol = []string{
	"tell me",
	"set master",
	"del master",
	"set nick",
	"set password",
	"join",
	"leave",
	"shutdown",
	"say",
	"say when",
	"cmd say",
	"exec",
	"exec when",
	"cmd exec",
	"sys run",
}

type Command struct {
	Method string
	Params []string
	Dest   string
}

// ParseCommand receives a command in a raw string and returns a new Command struct
func (a *Amigo) ParseCommand(msg *irc.Message) (*Command, error) {
	c := new(Command)

	// Remove nick from message
	raw := strings.TrimSpace(msg.Trailing[len(a.mem.Nick):])

	// Command
	found := false
	for _, cmd := range protocol {
		if strings.HasPrefix(raw, cmd) {
			found = true
			c.Method = strings.ToLower(cmd)
			break
		}
	}
	if !found {
		return nil, errors.New("Command not found")
	}

	// Params
	raw = strings.TrimSpace(raw[len(c.Method):])

	c.Params = strings.Split(raw, " "+param_delimiter)
	for key, param := range c.Params {
		c.Params[key] = strings.TrimSpace(strings.Replace(param, "\\"+param_delimiter, param_delimiter, -1))
	}

	// Talking on a channel or private?
	receiver := strings.Join(msg.Params, " ")

	if receiver == a.mem.Nick && msg.Prefix != nil {
		c.Dest = msg.Prefix.Name
	} else {
		// Talking on a channel
		c.Dest = receiver
	}

	// Check auth
	if msg.Prefix == nil || msg.Prefix.Name == "" {
		return nil, errors.New("User not authorized.")
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
	if a.mem.Password == c.Params[len(c.Params)-1] {
		auth = true
		c.Params = c.Params[:len(c.Params)-1]
	}
	if !auth {
		return nil, errors.New("User not authorized.")
	}

	return c, nil
}
