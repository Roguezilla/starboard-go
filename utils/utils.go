package utils

import "github.com/bwmarrin/discordgo"

func EmojiCount(s *discordgo.Session, m *discordgo.MessageReactionAdd) (int, error) {
	// test change
	users, err := s.MessageReactions(m.ChannelID, m.MessageID, m.Emoji.APIName(), 100, "", "")
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
