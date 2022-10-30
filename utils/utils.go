package utils

import "github.com/bwmarrin/discordgo"

func FormattedEmoji(emoji discordgo.Emoji) string {
	if emoji.ID == "" {
		return emoji.Name
	} else {
		temp := "<"
		if emoji.Animated {
			temp += "a"
		}
		return temp + ":" + emoji.Name + ":" + emoji.ID + ">"
	}
}

func EmojiCount(s *discordgo.Session, reactionEvent *discordgo.MessageReactionAdd) (int, error) {
	var editedEmoji = reactionEvent.Emoji.Name
	if reactionEvent.Emoji.ID != "" {
		editedEmoji = editedEmoji + ":" + reactionEvent.Emoji.ID
		if reactionEvent.Emoji.Animated {
			editedEmoji = "a:" + editedEmoji
		}
	}

	users, err := s.MessageReactions(reactionEvent.ChannelID, reactionEvent.MessageID, editedEmoji, 100, "", "")
	if err != nil {
		return -1, err
	}

	return len(users), nil
}

func CheckPermission(s *discordgo.Session, message *discordgo.Message, permission int64) (bool, error) {
	perms, err := s.UserChannelPermissions(message.Author.ID, message.ChannelID)
	if err == nil && (perms&permission == permission) {
		return true, nil
	}

	return false, err
}
