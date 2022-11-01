package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"roguezilla.github.io/starboard/cogs"
	"roguezilla.github.io/starboard/commands"
	"roguezilla.github.io/starboard/sqldb"
	"roguezilla.github.io/starboard/utils"
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

	commands.Handler(s, m)
	cogs.PixivHandler(s, m)
	cogs.TwitterHandler(s, m)
}

func messageReactionAdd(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
	if setup, err := sqldb.IsSetup(m.GuildID); err != nil || (err == nil && !setup) {
		return
	}
	// eventually custom embed stuff for reddit and instagram

	// starboard logic
	archived, err := sqldb.IsArchived(m.GuildID, m.ChannelID, m.MessageID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}
	if archived {
		return
	}

	emoji, err := sqldb.Emoji(m.GuildID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}
	if emoji != utils.FormattedEmoji(m.Emoji) {
		return
	}

	msg, err := s.ChannelMessage(m.ChannelID, m.MessageID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}
	msg.GuildID = m.GuildID // ChannelMessage returns Message with empty GuildID

	amount, err := sqldb.ChannelAmount(m.GuildID, m.ChannelID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	if amount == -1 {
		amount, err = sqldb.Amount(m.GuildID)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, err.Error())
			return
		}
	}

	emojiCount, err := utils.EmojiCount(s, m)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	if emojiCount >= amount {
		channelID, err := sqldb.Channel(m.GuildID)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, err.Error())
			return
		}

		cogs.Archive(s, msg, channelID)
	}
}
