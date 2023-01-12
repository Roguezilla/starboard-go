package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"starboard/events"
	"starboard/sqldb"

	"github.com/bwmarrin/discordgo"
)

func main() {
	sqldb.Connect("db.db")
	defer sqldb.Close()

	token, err := sqldb.Token()
	if err != nil {
		if strings.Contains(err.Error(), "no such table") {
			fmt.Printf("Creating tables...")
			if err = sqldb.CreateTables(); err != nil {
				log.Println("main.go main sqldb.CreateTables:", err)
				return
			}
			fmt.Println("✅")

			var input string
			fmt.Printf("Bot token: ")
			fmt.Scanln(&input)

			if err = sqldb.SetToken(input); err != nil {
				log.Println("main.go main sqldb.SetToken:", err)
				return
			}

			fmt.Println("✅")

			token, _ = sqldb.Token()
		} else {
			log.Println("main.go main sqldb.GetToken::", err)
			return
		}
	}

	session, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Println("main.go main discordgo.New::", err)
		return
	}

	// handlers
	session.AddHandlerOnce(events.OnReady)
	session.AddHandler(events.MessageCreate)
	session.AddHandler(events.MessageReactionAdd)

	// intents
	session.Identify.Intents |= discordgo.IntentMessageContent
	session.Identify.Intents |= discordgo.IntentsGuildMembers
	session.Identify.Intents |= discordgo.IntentGuildMessageReactions

	err = session.Open()
	if err != nil {
		log.Println("main.go main session.Open::", err)
		return
	}
	defer session.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
}
