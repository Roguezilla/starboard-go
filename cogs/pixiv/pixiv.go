package pixiv

import (
	"regexp"

	"github.com/bwmarrin/discordgo"
)

var regex, _ = regexp.Compile(`^https?://(?:www\.)?pixiv\.net/(?:en/)?artworks/(\d+)$`)

func HandleMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	pixivData := regex.FindStringSubmatch(m.Content)
	if len(pixivData) == 0 {
		return
	}

	s.ChannelMessageSend(m.ChannelID, "https://pixiv.kmn5.li/"+pixivData[1]+"?u="+m.Author.ID)
	s.ChannelMessageDelete(m.ChannelID, m.ID)
}
