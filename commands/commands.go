package commands

import (
	"os/exec"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"roguezilla.github.io/starboard/cogs/starboard"
	"roguezilla.github.io/starboard/sqldb"
	"roguezilla.github.io/starboard/utils"
)

func HandleMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	switch split := strings.Split(m.Content, " "); split[0] {
	case "sb!source":
		source(s, m, 0, split[1:]...)
	case "sb!setup":
		setup(s, m, 3, split[1:]...)
	case "sb!delete":
		deleteArchiveEntry(s, m, 0, split[1:]...)
	case "sb!set_emoji":
		setEmoji(s, m, 1, split[1:]...)
	case "sb!set_channel":
		setChannel(s, m, 1, split[1:]...)
	case "sb!set_amount":
		setAmount(s, m, 1, split[1:]...)
	case "sb!set_channel_amount":
		setCustomAmount(s, m, 2, split[1:]...)
	case "sb!override":
		starboard.ArchiveOverrideCommand(s, m, 1, split[1:]...)
	case "sb!pull":
		pull(s, m, 0, split[1:]...)
	}
}

func source(s *discordgo.Session, m *discordgo.MessageCreate, numArgs int, args ...string) {
	if numArgs != len(args) {
		s.ChannelMessageSendReply(m.ChannelID, "❌Invalid number of arguments, got "+strconv.Itoa(len(args))+" expected "+strconv.Itoa(numArgs)+".", m.Message.Reference())
		return
	}

	s.ChannelMessageSendReply(m.ChannelID, "<https://github.com/Roguezilla/starboard>", m.Message.Reference())
}

func setup(s *discordgo.Session, m *discordgo.MessageCreate, numArgs int, args ...string) {
	if numArgs != len(args) {
		s.ChannelMessageSend(m.ChannelID, "❌Invalid number of arguments, got "+strconv.Itoa(len(args))+" expected "+strconv.Itoa(numArgs)+".")
		return
	}

	if has, err := utils.CheckPermission(s, m.Message, discordgo.PermissionManageMessages); !has {
		s.ChannelMessageSendReply(m.ChannelID, "❌You don't have permission to do that.", m.Message.Reference())
		return
	} else if err != nil {
		s.ChannelMessageSendReply(m.ChannelID, err.Error(), m.Message.Reference())
		return
	}

	setup, err := sqldb.IsSetup(m.GuildID)
	if err != nil {
		s.ChannelMessageSendReply(m.ChannelID, err.Error(), m.Message.Reference())
		return
	} else if setup {
		s.ChannelMessageSendReply(m.ChannelID, "❌Server already set-up.", m.Message.Reference())
		return
	}
	/* */
	parsed, err := strconv.ParseInt(args[2], 10, 0)
	if err != nil {
		s.ChannelMessageSendReply(m.ChannelID, "❌Non-numeric amount.", m.Message.Reference())
		return
	}

	// passed "raw" emoji into APIName emoji
	arr := [5]string{"<", "a", ":", ">"}
	for i := 0; i < len(arr); i++ {
		args[1] = strings.Replace(args[1], arr[i], "", 1)
	}

	if err := sqldb.Setup(m.GuildID, args[0][2:len(args[0])-1], args[1], parsed); err == nil {
		s.ChannelMessageSendReply(m.ChannelID, "*✅*", m.Message.Reference())
	}
}

func deleteArchiveEntry(s *discordgo.Session, m *discordgo.MessageCreate, numArgs int, args ...string) {
	if !commandRunnable(s, m, numArgs, args...) {
		return
	}

	if m.MessageReference == nil {
		s.ChannelMessageSendReply(m.ChannelID, "❗You have to reply to the message you want to unarchive.", m.Message.Reference())
		return
	}

	if err := sqldb.Unarchive(m.MessageReference.GuildID, m.MessageReference.ChannelID, m.MessageReference.MessageID); err == nil {
		s.ChannelMessageSendReply(m.ChannelID, "*✅*", m.Message.Reference())
	}
}

func setEmoji(s *discordgo.Session, m *discordgo.MessageCreate, numArgs int, args ...string) {
	if !commandRunnable(s, m, numArgs, args...) {
		return
	}

	if err := sqldb.SetEmoji(m.GuildID, args[0]); err == nil {
		s.ChannelMessageSendReply(m.ChannelID, "*✅*", m.Message.Reference())
	}
}

func setChannel(s *discordgo.Session, m *discordgo.MessageCreate, numArgs int, args ...string) {
	if !commandRunnable(s, m, numArgs, args...) {
		return
	}
	/* */
	if err := sqldb.SetChannel(m.GuildID, args[0]); err == nil {
		s.ChannelMessageSendReply(m.ChannelID, "*✅*", m.Message.Reference())
	}
}

func setAmount(s *discordgo.Session, m *discordgo.MessageCreate, numArgs int, args ...string) {
	if !commandRunnable(s, m, numArgs, args...) {
		return
	}

	parsed, err := strconv.ParseInt(args[0], 10, 0)
	if err != nil {
		s.ChannelMessageSendReply(m.ChannelID, "❗Non-numeric amount.", m.Message.Reference())
		return
	}

	if err := sqldb.SetAmount(m.GuildID, parsed); err == nil {
		s.ChannelMessageSendReply(m.ChannelID, "*✅*", m.Message.Reference())
	}
}

func setCustomAmount(s *discordgo.Session, m *discordgo.MessageCreate, numArgs int, args ...string) {
	if !commandRunnable(s, m, numArgs, args...) {
		return
	}

	parsed, err := strconv.ParseInt(args[1], 10, 0)
	if err != nil {
		s.ChannelMessageSendReply(m.ChannelID, "❗Non-numeric amount.", m.Message.Reference())
		return
	}

	if err := sqldb.SetChannelAmount(m.GuildID, args[0][2:len(args[0])-1], parsed); err == nil {
		s.ChannelMessageSendReply(m.ChannelID, "*✅*", m.Message.Reference())
	}
}

func pull(s *discordgo.Session, m *discordgo.MessageCreate, numArgs int, args ...string) {
	if !commandRunnable(s, m, numArgs, args...) {
		return
	}

	cmd := exec.Command("git", "pull")
	cmdOutput, err := cmd.Output()
	if err != nil {
		s.ChannelMessageSendReply(m.ChannelID, err.Error(), m.Message.Reference())
		return
	}

	embed := discordgo.MessageEmbed{
		Color: 0xffcc00,
		Footer: &discordgo.MessageEmbedFooter{
			Text: "by rogue#0001",
		},
		Title:       "Update Log",
		Description: "```" + string(cmdOutput)[0:len(string(cmdOutput))-6] + "```",
	}

	if _, err := s.ChannelMessageSendEmbed(m.ChannelID, &embed); err != nil {
		s.ChannelMessageSendReply(m.ChannelID, err.Error(), m.Reference())
		return
	}
}

/* */
func commandRunnable(s *discordgo.Session, m *discordgo.MessageCreate, numArgs int, args ...string) bool {
	if numArgs != len(args) {
		s.ChannelMessageSendReply(m.ChannelID, "❌Invalid number of arguments, got "+strconv.Itoa(len(args))+" expected "+strconv.Itoa(numArgs)+".", m.Message.Reference())
		return false
	}

	if has, err := utils.CheckPermission(s, m.Message, discordgo.PermissionManageMessages); !has {
		s.ChannelMessageSendReply(m.ChannelID, "❌You don't have permission to do that.", m.Message.Reference())
		return false
	} else if err != nil {
		s.ChannelMessageSendReply(m.ChannelID, err.Error(), m.Message.Reference())
		return false
	}

	setup, err := sqldb.IsSetup(m.GuildID)
	if err != nil {
		s.ChannelMessageSendReply(m.ChannelID, err.Error(), &discordgo.MessageReference{GuildID: m.GuildID, ChannelID: m.ChannelID, MessageID: m.ID})
		return false
	} else if !setup {
		s.ChannelMessageSendReply(m.ChannelID, "❌Server has not been set-up.", m.Message.Reference())
		return false
	}

	return true
}
