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

func EmojiCount(s *discordgo.Session, channelId string, messageId string, emoji discordgo.Emoji) int {
	var editedEmoji = emoji.Name
	if emoji.ID != "" {
		editedEmoji = editedEmoji + ":" + emoji.ID
		if emoji.Animated {
			editedEmoji = "a:" + editedEmoji
		}
	}

	if users, err := s.MessageReactions(channelId, messageId, editedEmoji, 100, "", ""); err != nil {
		return -1
	} else {
		return len(users)
	}
}

func CheckPermission(s *discordgo.Session, message *discordgo.Message, permission int64) (bool, error) {
	perms, err := s.UserChannelPermissions(message.Author.ID, message.ChannelID)
	if err == nil && (perms&permission == permission) {
		return true, nil
	}

	return false, err
}
