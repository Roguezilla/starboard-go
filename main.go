package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"starboard/sqldb"

	"github.com/bwmarrin/discordgo"
)

func main() {
	err := sqldb.Open("db.db")
	if err != nil {
		fmt.Println("sqldb.Open: ", err)
		return
	}
	defer sqldb.Close()

	token, err := sqldb.Token()
	if err != nil {
		if err.Error() == "no such table: settings" {
			fmt.Printf("Creating tables...")
			err := sqldb.CreateTables()
			if err != nil {
				fmt.Println("sqldb.CreateTables: ", err)
				return
			}
			fmt.Println("✅")

			var input string
			fmt.Printf("Bot token: ")
			fmt.Scanln(&input)
			sqldb.SetToken(input)

			token, _ = sqldb.Token()
		} else {
			fmt.Println("sqldb.GetToken: ", err)
			return
		}
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
