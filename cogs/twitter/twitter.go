package twitter

import (
	"regexp"

	"github.com/bwmarrin/discordgo"
)

var regex, _ = regexp.Compile(`^https?://(?:mobile.)?twitter\.com(/.+/status/\d+)$`)

func HandleMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	path := regex.FindStringSubmatch(m.Content)
	if len(path) == 0 {
		return
	}

	s.ChannelMessageSend(m.ChannelID, "https://vxtwitter.com"+path[1]+"?u="+m.Author.ID)
	s.ChannelMessageDelete(m.ChannelID, m.ID)
}
