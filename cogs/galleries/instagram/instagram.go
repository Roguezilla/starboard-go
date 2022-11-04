package instagram

import (
	"regexp"

	"github.com/bwmarrin/discordgo"
)

type gallery struct {
	GallerySize int
	CurrentIdx  int
	Gallery     []string
}

var instagramCache = map[string]gallery{}

var regex, _ = regexp.Compile(`^((?:(?:(?:https):(?://)+)(?:www\.)?)(?:instagram)\.com/p/.+)$`)

func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	match := regex.FindStringSubmatch(m.Content)
	if len(match) == 0 {
		return
	}
}

func MessageReactionAdd(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
	// right and left array emojis
	if m.Emoji.APIName() != string([]rune{10145, 65039}) && m.Emoji.APIName() != string([]rune{11013, 65039}) {
		return
	}
}
