package twitter

import (
	"regexp"
	"time"

	"github.com/bwmarrin/discordgo"
)

var regex, _ = regexp.Compile(`^https://(?:mobile.)?twitter\.com(/.+/status/\d+)$`)

func HandleMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	twitterData := regex.FindStringSubmatch(m.Content)
	if len(twitterData) == 0 {
		return
	}

	// discord is extremely slow with uncached video tweets, because of this the Message
	// object sent to the bot might not have the Video field filled, which is why it's
	// preferable to wait 3 seconds and "refresh" the message
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
