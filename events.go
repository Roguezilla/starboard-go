package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"roguezilla.github.io/starboard/sqldb"
)

func onReady(s *discordgo.Session, m *discordgo.Ready) {
	s.UpdateStatusComplex(discordgo.UpdateStatusData{
		Status: "idle",
		Activities: []*discordgo.Activity{{
			Name: "the stars.",
			Type: discordgo.ActivityTypeWatching,
			URL:  "",
		}},
	})
	fmt.Println("->" + m.User.Username + " is ready.")
}

// TODO: command framework:tm:
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if setup, err := sqldb.IsSetup(m.GuildID); err != nil || (err == nil && !setup) {
		return
	}

}

func messageReactionAdd(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
	if setup, err := sqldb.IsSetup(m.GuildID); err != nil || (err == nil && !setup) {
		return
	}

	/*
		if archived, err := sqldb.IsArchived(m.GuildID, m.ChannelID, m.MessageID); err == nil && !archived {
			sqldb.Archive(m.GuildID, m.ChannelID, m.MessageID)
		}
	*/
}
