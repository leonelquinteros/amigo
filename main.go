// Amigo is an IRC bot that learns to talk from its masters.
// The entire bot can be set up through IRC commands as follows:
//
// Command format has to be in the following form to be recognized:
//
//  [ NICK ] [ COMMAND ] [ PARAM ] [ SPACE + PARAM_DELIMITER [ PARAM ] ... ]
//
//  Where NICK is the bot's nick on the IRC server,
//  COMMAND is a string present in the protocol variable defined on this file,
//  [ PARAM ] is any string formed by any chars other than a combination of SPACE and PARAM_DELIMITER,
//  SPACE is the space character and PARAM_DELIMITER is defined in param_delimiter constant on command.go, a double semicolon ";;" by default.
//
//  If PARAM_DELIMITER has to be included as part of a param value (usually not needed), you can escape it using a backslash \,
//  Any combination of backslash and the param delimiter will be replaced by the param delimiter itself at parsing time.
//
package main

import (
	"flag"
)

func main() {
	// Load Config
	host := flag.String("h", "", "Host: The IRC host to connect to")
	channel := flag.String("c", "", "Channel: The IRC #channel to join after connect to the IRC server")
	nick := flag.String("n", "", "Nick: The nick to use on the IRC server")
	password := flag.String("p", "", "Password: The master password to perform authenticated commands without beign a registered master of the bot")
	flag.Parse()

	if *host == "" || *channel == "" || *nick == "" || *password == "" {
		println("Usage:")
		println("amigo -h=irc.host.com:6667 -c=#channel_name -n=nick -p=masterpassword")
		return
	}

	amigo := new(Amigo)
	amigo.EhAmigo(*host, *channel, *nick, *password)
}
