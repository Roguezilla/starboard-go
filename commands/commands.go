package commands

import (
	"database/sql"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"roguezilla.github.io/starboard/sqldb"
	"roguezilla.github.io/starboard/utils"
)

func Handler(s *discordgo.Session, m *discordgo.MessageCreate) {
	switch split := strings.Split(m.Content, " "); split[0] {
	case "sb!source":
		source(s, m, 0, split[1:]...)
	case "sb!setup":
		setup(s, m, 3, split[1:]...)
	}
}

func source(s *discordgo.Session, m *discordgo.MessageCreate, numArgs int, args ...string) {
	if numArgs != len(args) {
		s.ChannelMessageSend(m.ChannelID, "Invalid number of arguments, got "+strconv.Itoa(len(args))+" expected "+strconv.Itoa(numArgs)+".")
		return
	}

	s.ChannelMessageSend(m.ChannelID, "<https://github.com/Roguezilla/starboard>")
}

func setup(s *discordgo.Session, m *discordgo.MessageCreate, numArgs int, args ...string) {
	if numArgs != len(args) {
		s.ChannelMessageSend(m.ChannelID, "Invalid number of arguments, got "+strconv.Itoa(len(args))+" expected "+strconv.Itoa(numArgs)+".")
		return
	}

	if has, err := utils.CheckPermission(s, m.Message, discordgo.PermissionManageMessages); !has {
		return
	} else if err != nil {
		s.ChannelMessageSend(m.ChannelID, err.Error())
	}

	if setup, err := sqldb.IsSetup(m.GuildID); setup {
		s.ChannelMessageSend(m.ChannelID, "Server already setup.")
		return
	} else if err != nil && err != sql.ErrNoRows {
		s.ChannelMessageSend(m.ChannelID, err.Error())
	}

	parsed, err := strconv.ParseInt(args[2], 10, 0)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Non-numeric amount.")
		return
	}

	sqldb.Setup(m.GuildID, args[0][2:len(args[0])-1], args[1], parsed)
	s.ChannelMessageSend(m.ChannelID, "Successful.")
}
