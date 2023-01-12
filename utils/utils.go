package utils

import "github.com/bwmarrin/discordgo"

func EmojiCount(s *discordgo.Session, m *discordgo.MessageReactionAdd) (int, error) {
	users, err := s.MessageReactions(m.ChannelID, m.MessageID, m.Emoji.APIName(), 100, "", "")

	return len(users), err
}

func CheckPermission(s *discordgo.Session, message *discordgo.Message, permission int64) (bool, error) {
	perms, err := s.UserChannelPermissions(message.Author.ID, message.ChannelID)

	return perms&permission == permission, err
}
