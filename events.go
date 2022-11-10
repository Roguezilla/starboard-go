package main

import (
	"fmt"

	"starboard/cogs/galleries"
	"starboard/cogs/pixiv"
	"starboard/cogs/starboard"
	"starboard/cogs/twitter"
	"starboard/commands"
	"starboard/sqldb"

	"github.com/bwmarrin/discordgo"
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

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	commands.HandleMessageCreate(s, m)
	pixiv.HandleMessageCreate(s, m)
	twitter.HandleMessageCreate(s, m)
	galleries.HandleMessageCreate(s, m)
}

func messageReactionAdd(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
	if m.Member.User.ID == s.State.User.ID {
		return
	}

	setup, err := sqldb.IsSetup(m.GuildID)
	if err != nil {
		s.ChannelMessageSendReply(m.ChannelID, err.Error(), &discordgo.MessageReference{GuildID: m.GuildID, ChannelID: m.ChannelID, MessageID: m.MessageID})
		return
	} else if !setup {
		return
	}

	galleries.HandleMessageReactionAdd(s, m)
	starboard.HandleMessageReactionAdd(s, m)
}
