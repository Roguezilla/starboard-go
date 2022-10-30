package cogs

import (
	"github.com/bwmarrin/discordgo"
	"roguezilla.github.io/starboard/sqldb"
)

type embedInfo struct {
	Flag         string
	Content      string
	MediaURL     string
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
		Footer: &discordgo.MessageEmbedFooter{
			Text: "by rogue#0001",
		},
	}

	if embedInfo.CustomAuthor != nil {
		embed.Author = &discordgo.MessageEmbedAuthor{
			Name:    embedInfo.CustomAuthor.Username,
			IconURL: embedInfo.CustomAuthor.AvatarURL(""),
		}
	} else {
		embed.Author = &discordgo.MessageEmbedAuthor{
			Name:    m.Author.Username,
			IconURL: m.Author.AvatarURL(""),
		}
	}
	if embedInfo.Content != "" {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "What?",
			Value:  embedInfo.Content,
			Inline: false,
		})
	}
	if embedInfo.Flag == "image" && embedInfo.MediaURL != "" {
		embed.Image = &discordgo.MessageEmbedImage{
			URL: embedInfo.MediaURL,
		}
	}

	embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
		Name:   "Where?",
		Value:  "[Jump to ](https://discordapp.com/channels/" + m.GuildID + "/" + m.ChannelID + "/" + m.ID + ")<#" + m.ChannelID + ">",
		Inline: false,
	})

	if err := sqldb.Archive(m.GuildID, m.ChannelID, m.ID); err != nil {
		s.ChannelMessageSendReply(m.ChannelID, err.Error(), m.Reference())
		return
	}

	if _, err := s.ChannelMessageSendEmbed(channelID, &embed); err != nil {
		s.ChannelMessageSendReply(m.ChannelID, err.Error(), m.Reference())
		return
	}

	if embedInfo.Flag == "video" && embedInfo.MediaURL != "" {
		if _, err := s.ChannelMessageSend(channelID, embedInfo.MediaURL); err != nil {
			s.ChannelMessageSendReply(m.ChannelID, err.Error(), m.Reference())
			return
		}
	}
}
