package cogs

import (
	"regexp"

	"github.com/bwmarrin/discordgo"
)

var pixivRegex, _ = regexp.Compile(`^https:\/\/(?:www\.)?pixiv\.net\/(?:en\/)?artworks\/(\d+)$`)

func PixivHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	pixivData := pixivRegex.FindStringSubmatch(m.Content)
	if len(pixivData) > 0 {
		s.ChannelMessageSend(m.ChannelID, "https://pixiv.kmn5.li/"+pixivData[1]+"?u="+m.Author.ID)
		s.ChannelMessageDelete(m.ChannelID, m.ID)
	}
}
