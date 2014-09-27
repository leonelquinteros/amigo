package main

func main() {
    amigo := new(Amigo)

    amigo.host      = "irc.freenode.org:6667"
    amigo.channel   = "#amigo-bot"
    amigo.nick      = "amigobot"

    amigo.EhAmigo()
}
