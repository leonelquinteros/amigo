// Amigo is an IRC bot that learns to talk from its masters.
// The entire bot can be set up through IRC commands as follows:
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
