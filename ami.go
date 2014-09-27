package main

import (
    "fmt"
    "errors"
    "github.com/sorcix/irc"
)

// This is the bot
type Amigo struct {
    // Connection params
    host, channel, nick string

    // Connection handler
    conn *irc.Conn
}

// Bot starter
func (self *Amigo) EhAmigo() {
    err := self.Connect()

    if err != nil {
        fmt.Println("AMIGO ABORT!")

        return
    }

    self.Listen()
}


// IRC connect
func (self *Amigo) Connect() (e error) {
    c, err := irc.Dial(self.host)

    if err != nil {
        errMsg := "AMIGO ERROR: " + err.Error()
        fmt.Println(errMsg)
        return errors.New(errMsg)
    }

    self.conn = c

    return nil
}


// Messages listener. Gets all the network stream and dispatches the messages.
func (self *Amigo) Listen() {
    for {
        msg, err := self.conn.Decode()

        if err != nil {
            fmt.Println("AMIGO ERROR: " + err.Error())
            self.conn.Close()
            break
        }

        fmt.Println(msg)
    }
}
