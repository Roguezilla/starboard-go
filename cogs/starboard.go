package cogs

import (
	"github.com/bwmarrin/discordgo"
	"roguezilla.github.io/starboard/sqldb"
)

type embedInfo struct {
	Flag         string
	Content      string
	ImageUrl     string
	CustomAuthor *discordgo.User
}

func buildEmbedInfo(m *discordgo.Message) embedInfo {
	e := embedInfo{
		Flag:    "message",
		Content: m.Content,
	}

	return e
}

func Archive(s *discordgo.Session, m *discordgo.Message, channelID string) {
	embedInfo := buildEmbedInfo(m)

	embed := discordgo.MessageEmbed{
		Color: 0xffcc00,
	}

	embed.Author = &discordgo.MessageEmbedAuthor{
		Name:    m.Author.Username,
		IconURL: m.Author.AvatarURL("128"),
	}
	embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
		Name:   "What?",
		Value:  embedInfo.Content,
		Inline: false,
	})
	embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
		Name:   "Where?",
		Value:  "[Jump to ](https://discordapp.com/channels/" + m.GuildID + "/" + m.ChannelID + "/" + m.ID + ")<#" + m.ChannelID + ">",
		Inline: false,
	})
	embed.Footer = &discordgo.MessageEmbedFooter{
		Text: "by rogue#0001",
	}

	if _, err := s.ChannelMessageSendEmbed(channelID, &embed); err != nil {
		s.ChannelMessageSendReply(m.ChannelID, err.Error(), m.Reference())
		return
	}

	if err := sqldb.Archive(m.GuildID, m.ChannelID, m.ID); err != nil {
		s.ChannelMessageSendReply(m.ChannelID, err.Error(), m.Reference())
	}
}
