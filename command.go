package main

import (
    "strings"
    "errors"
    "github.com/sorcix/irc"
)


var param_delimiter = ";;"

// Protocol commands list
var protocol = []string{
        "set master",
        "del master",
        "set nick",
        "set password",
        "join",
        "leave",
        "sys run",
        "say",
        "say when",
        "cmd say",
        "exec",
        "exec when",
        "cmd exec",
}


type Command struct {
    Method      string
    Params      []string
    Dest        string
}

// ParseCommand receives a command in a raw string and returns a new Command struct
// Command format has to be in the following form to be recognized:
//
// NICK [ SPACE... ] PROTOCOL_COMMAND [ SPACE... ] 'param' [ SPACE PARAM_DELIMITER 'extra params' ... ]
//
// Where NICK is the bot's nick on the IRC server,
// PROTOCOL_COMMAND is a string present in the protocol variable defined on this file,
// 'param' and 'extra params' are any string formed by any chars other than a combination of SPACE and PARAM_DELIMITER,
// SPACE is the space character and PARAM_DELIMITER is defined in param_delimiter variable on this file, usually a double semicolon ';;'.
//
// If PARAM_DELIMITER has to be included as part of a param value (usually not needed), you can escape it using a backslash \,
// Any combination of backslash and the param delimiter will be replaced by the param delimiter itself at parsing time.
func (a *Amigo) ParseCommand(msg *irc.Message) (*Command, error) {
    c := new(Command)

    // Remove nick from message
    raw := strings.TrimSpace(msg.Trailing[len(a.mem.Nick):])

    // Command
    found := false
    for _, cmd := range protocol {
        if strings.HasPrefix(raw, cmd) {
            found = true
            c.Method = cmd
            break
        }
    }
    if !found {
        return nil, errors.New("Command not found")
    }

    // Params
    raw = strings.TrimSpace(raw[len(c.Method):])

    c.Params = strings.Split(raw, " " + param_delimiter)
    for key, param := range c.Params {
        c.Params[key] = strings.TrimSpace(strings.Replace(param, "\\" + param_delimiter, param_delimiter, -1))
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
