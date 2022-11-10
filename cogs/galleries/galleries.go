package galleries

import (
	"starboard/cogs/galleries/instagram"
	"starboard/cogs/galleries/reddit"

	"github.com/bwmarrin/discordgo"
)

func HandleMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	reddit.MessageCreate(s, m)
	instagram.MessageCreate(s, m)
}

func HandleMessageReactionAdd(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
	reddit.MessageReactionAdd(s, m)
	instagram.MessageReactionAdd(s, m)
}
