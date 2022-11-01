package cogs

import (
	"regexp"
	"time"

	"github.com/bwmarrin/discordgo"
)

var twitterRegex2, _ = regexp.Compile(`^https://(?:mobile.)?twitter\.com(/.+/status/\d+)$`)

func TwitterHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	twitterData := twitterRegex2.FindStringSubmatch(m.Content)
	if len(twitterData) > 0 {
		time.Sleep(3 * time.Second)

		refreshedMessage, err := s.ChannelMessage(m.ChannelID, m.ID)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, err.Error())
			return
		}
		if len(refreshedMessage.Embeds) < 1 || refreshedMessage.Embeds[0].Video == nil {
			return
		}

		s.ChannelMessageSend(m.ChannelID, "https://vxtwitter.com"+twitterData[1]+"?u="+m.Author.ID)
		s.ChannelMessageDelete(m.ChannelID, m.ID)
	}

}
