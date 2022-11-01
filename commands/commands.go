package commands

import (
	"database/sql"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"roguezilla.github.io/starboard/cogs"
	"roguezilla.github.io/starboard/sqldb"
	"roguezilla.github.io/starboard/utils"
)

func Handler(s *discordgo.Session, m *discordgo.MessageCreate) {
	switch split := strings.Split(m.Content, " "); split[0] {
	case "sb!source":
		source(s, m, 0, split[1:]...)
	case "sb!setup":
		setup(s, m, 3, split[1:]...)
	case "sb!delete":
		delete(s, m, 0, split[1:]...)
	case "sb!set_emoji":
		setEmoji(s, m, 1, split[1:]...)
	case "sb!set_channel":
		setChannel(s, m, 1, split[1:]...)
	case "sb!set_amount":
		setAmount(s, m, 1, split[1:]...)
	case "sb!set_channel_amount":
		setCustomAmount(s, m, 2, split[1:]...)
	case "sb!override":
		cogs.ArchiveOverrideCommand(s, m, 1, split[1:]...)
	}
}

func source(s *discordgo.Session, m *discordgo.MessageCreate, numArgs int, args ...string) {
	if numArgs != len(args) {
		s.ChannelMessageSendReply(m.ChannelID, "Invalid number of arguments, got "+strconv.Itoa(len(args))+" expected "+strconv.Itoa(numArgs)+".", m.Message.Reference())
		return
	}
	/* */
	s.ChannelMessageSendReply(m.ChannelID, "<https://github.com/Roguezilla/starboard>", m.Message.Reference())
}

func setup(s *discordgo.Session, m *discordgo.MessageCreate, numArgs int, args ...string) {
	if numArgs != len(args) {
		s.ChannelMessageSendReply(m.ChannelID, "Invalid number of arguments, got "+strconv.Itoa(len(args))+" expected "+strconv.Itoa(numArgs)+".", m.Message.Reference())
		return
	}

	if has, err := utils.CheckPermission(s, m.Message, discordgo.PermissionManageMessages); !has {
		s.ChannelMessageSendReply(m.ChannelID, "You don't have permission to do that.", m.Message.Reference())
		return
	} else if err != nil {
		s.ChannelMessageSendReply(m.ChannelID, err.Error(), m.Message.Reference())
		return
	}
	/* */
	if setup, err := sqldb.IsSetup(m.GuildID); setup {
		s.ChannelMessageSendReply(m.ChannelID, "Server already setup.", m.Message.Reference())
		return
	} else if err != nil && err != sql.ErrNoRows {
		s.ChannelMessageSendReply(m.ChannelID, err.Error(), m.Message.Reference())
		return
	}

	parsed, err := strconv.ParseInt(args[2], 10, 0)
	if err != nil {
		s.ChannelMessageSendReply(m.ChannelID, "Non-numeric amount.", m.Message.Reference())
		return
	}

	if err := sqldb.Setup(m.GuildID, args[0][2:len(args[0])-1], args[1], parsed); err == nil {
		s.ChannelMessageSendReply(m.ChannelID, "*✅*", m.Message.Reference())
	}
}

func delete(s *discordgo.Session, m *discordgo.MessageCreate, numArgs int, args ...string) {
	if numArgs != len(args) {
		s.ChannelMessageSend(m.ChannelID, "Invalid number of arguments, got "+strconv.Itoa(len(args))+" expected "+strconv.Itoa(numArgs)+".")
		return
	}

	if has, err := utils.CheckPermission(s, m.Message, discordgo.PermissionManageMessages); !has {
		s.ChannelMessageSendReply(m.ChannelID, "You don't have permission to do that.", m.Message.Reference())
		return
	} else if err != nil {
		s.ChannelMessageSendReply(m.ChannelID, err.Error(), m.Message.Reference())
		return
	}
	/* */
	if m.MessageReference == nil {
		s.ChannelMessageSendReply(m.ChannelID, "You have to reply to the message you want to unarchive.", m.Message.Reference())
		return
	}

	if err := sqldb.Unarchive(m.MessageReference.GuildID, m.MessageReference.ChannelID, m.MessageReference.MessageID); err == nil {
		s.ChannelMessageSendReply(m.ChannelID, "*✅*", m.Message.Reference())
	}
}

func setEmoji(s *discordgo.Session, m *discordgo.MessageCreate, numArgs int, args ...string) {
	if numArgs != len(args) {
		s.ChannelMessageSend(m.ChannelID, "Invalid number of arguments, got "+strconv.Itoa(len(args))+" expected "+strconv.Itoa(numArgs)+".")
		return
	}

	if has, err := utils.CheckPermission(s, m.Message, discordgo.PermissionManageMessages); !has {
		s.ChannelMessageSendReply(m.ChannelID, "You don't have permission to do that.", m.Message.Reference())
		return
	} else if err != nil {
		s.ChannelMessageSendReply(m.ChannelID, err.Error(), m.Message.Reference())
		return
	}
	/* */
	if err := sqldb.SetEmoji(m.GuildID, args[0]); err == nil {
		s.ChannelMessageSendReply(m.ChannelID, "*✅*", m.Message.Reference())
	}
}

func setChannel(s *discordgo.Session, m *discordgo.MessageCreate, numArgs int, args ...string) {
	if numArgs != len(args) {
		s.ChannelMessageSend(m.ChannelID, "Invalid number of arguments, got "+strconv.Itoa(len(args))+" expected "+strconv.Itoa(numArgs)+".")
		return
	}

	if has, err := utils.CheckPermission(s, m.Message, discordgo.PermissionManageMessages); !has {
		s.ChannelMessageSendReply(m.ChannelID, "You don't have permission to do that.", m.Message.Reference())
		return
	} else if err != nil {
		s.ChannelMessageSendReply(m.ChannelID, err.Error(), m.Message.Reference())
		return
	}
	/* */
	if err := sqldb.SetChannel(m.GuildID, args[0]); err == nil {
		s.ChannelMessageSendReply(m.ChannelID, "*✅*", m.Message.Reference())
	}
}

func setAmount(s *discordgo.Session, m *discordgo.MessageCreate, numArgs int, args ...string) {
	if numArgs != len(args) {
		s.ChannelMessageSend(m.ChannelID, "Invalid number of arguments, got "+strconv.Itoa(len(args))+" expected "+strconv.Itoa(numArgs)+".")
		return
	}

	if has, err := utils.CheckPermission(s, m.Message, discordgo.PermissionManageMessages); !has {
		s.ChannelMessageSendReply(m.ChannelID, "You don't have permission to do that.", m.Message.Reference())
		return
	} else if err != nil {
		s.ChannelMessageSendReply(m.ChannelID, err.Error(), m.Message.Reference())
		return
	}
	/* */
	parsed, err := strconv.ParseInt(args[0], 10, 0)
	if err != nil {
		s.ChannelMessageSendReply(m.ChannelID, "Non-numeric amount.", m.Message.Reference())
		return
	}

	if err := sqldb.SetAmount(m.GuildID, parsed); err == nil {
		s.ChannelMessageSendReply(m.ChannelID, "*✅*", m.Message.Reference())
	}
}

func setCustomAmount(s *discordgo.Session, m *discordgo.MessageCreate, numArgs int, args ...string) {
	if numArgs != len(args) {
		s.ChannelMessageSend(m.ChannelID, "Invalid number of arguments, got "+strconv.Itoa(len(args))+" expected "+strconv.Itoa(numArgs)+".")
		return
	}

	if has, err := utils.CheckPermission(s, m.Message, discordgo.PermissionManageMessages); !has {
		s.ChannelMessageSendReply(m.ChannelID, "You don't have permission to do that.", m.Message.Reference())
		return
	} else if err != nil {
		s.ChannelMessageSendReply(m.ChannelID, err.Error(), m.Message.Reference())
		return
	}
	/* */
	parsed, err := strconv.ParseInt(args[1], 10, 0)
	if err != nil {
		s.ChannelMessageSendReply(m.ChannelID, "Non-numeric amount.", m.Message.Reference())
		return
	}

	if err := sqldb.SetChannelAmount(m.GuildID, args[0][2:len(args[0])-1], parsed); err == nil {
		s.ChannelMessageSendReply(m.ChannelID, "*✅*", m.Message.Reference())
	}
}
