package galleries

import (
	"github.com/bwmarrin/discordgo"
	"roguezilla.github.io/starboard/cogs/galleries/instagram"
	"roguezilla.github.io/starboard/cogs/galleries/reddit"
)

func HandleMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	reddit.MessageCreate(s, m)
	instagram.MessageCreate(s, m)
}

func HandleMessageReactionAdd(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
	reddit.MessageReactionAdd(s, m)
	instagram.MessageReactionAdd(s, m)
}
