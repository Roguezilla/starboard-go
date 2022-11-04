package reddit

import (
	"regexp"

	"github.com/bwmarrin/discordgo"
)

type gallery struct {
	GallerySize int
	CurrentIdx  int
	Gallery     []string
}

var redditCache = map[string]gallery{}

var regex, _ = regexp.Compile(`^((?:(?:(?:https):(?://)+)(?:www\.)?)redd(?:it\.com/|\.it/).+)$`)

func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	match := regex.FindStringSubmatch(m.Content)
	if len(match) == 0 {
		return
	}
}

func MessageReactionAdd(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
	// right and left array emojis in rune format
	if m.Emoji.APIName() != string([]rune{10145, 65039}) && m.Emoji.APIName() != string([]rune{11013, 65039}) {
		return
	}
}
