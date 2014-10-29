package main

func main() {
	amigo := new(Amigo)

	host    := "irc.freenode.org:6667"
	channel := "#amigo-bot"
	nick    := "amigobot"
    master  := "peiiion"

	amigo.EhAmigo(host, channel, nick, master)
}
