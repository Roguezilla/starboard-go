package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"roguezilla.github.io/starboard/sqldb"
)

func main() {
	err := sqldb.Open("db.db")
	if err != nil {
		fmt.Println("sqldb.Open: ", err)
		return
	}
	defer sqldb.Close()

	token, err := sqldb.GetToken()
	if err != nil {
		fmt.Println("sqldb.GetToken: ", err)
		return
	}

	session, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("discordgo.New: ", err)
		return
	}

	// handlers
	session.AddHandlerOnce(onReady)
	session.AddHandler(messageCreate)
	session.AddHandler(messageReactionAdd)

	// intents
	session.Identify.Intents |= discordgo.IntentMessageContent
	session.Identify.Intents |= discordgo.IntentsGuildMembers
	session.Identify.Intents |= discordgo.IntentGuildMessageReactions

	err = session.Open()
	if err != nil {
		fmt.Println("discordgo.Open: ", err)
		return
	}

	defer session.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
}
