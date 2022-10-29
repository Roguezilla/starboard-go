package sqldb

import (
	"database/sql"
	"errors"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

var conn *sql.DB

func Open(file string) error {
	db, err := sql.Open("sqlite3", file)
	if err == nil {
		conn = db
	}

	return err
}

func Close() {
	conn.Close()
}

func GetToken() (string, error) {
	rows, err := conn.Query("SELECT value FROM settings WHERE name LIKE 'token'")
	if err != nil {
		return "", err
	}
	defer rows.Close()

	var token string

	if rows.Next() {
		if err = rows.Scan(&token); err != nil {
			return "", err
		}
	} else {
		return "", errors.New("No token in database")
	}

	return token, nil
}

func IsSetup(guildId string) (bool, error) {
	_, err := GetEmoji(guildId)
	if err != nil {
		return false, err
	}

	return true, nil
}

func Setup(guildId string, channelId string, emoji string, amount int64) (bool, error) {
	parsedGuildId, err := strconv.ParseInt(guildId, 10, 0)
	if err != nil {
		return false, err
	}
	parsedChannelId, err := strconv.ParseInt(channelId, 10, 0)
	if err != nil {
		return false, err
	}

	if _, err := conn.Exec("INSERT INTO server (server_id, channel_id, archive_emote, archive_amount) VALUES (?, ?, ?, ?)", parsedGuildId, parsedChannelId, emoji, amount); err != nil {
		return false, err
	}

	return true, nil
}

func GetChannel(guildId string) (int, error) {
	stmt, err := conn.Prepare("SELECT archive_channel FROM server WHERE server_id = ?")
	if err != nil {
		return -1, err
	}

	parsed, err := strconv.ParseInt(guildId, 10, 0)
	if err != nil {
		return -1, err
	}

	var archiveChannel int
	if err = stmt.QueryRow(parsed).Scan(&archiveChannel); err != nil {
		return -1, err
	}

	return archiveChannel, nil
}

func SetChannel(guildId string, channel int) error {
	parsed, err := strconv.ParseInt(guildId, 10, 0)
	if err != nil {
		return err
	}

	if _, err := conn.Exec("UPDATE server SET archive_channel = ? WHERE server_id = ?", channel, parsed); err != nil {
		return err
	}

	return nil
}

func GetEmoji(guildId string) (string, error) {
	stmt, err := conn.Prepare("SELECT archive_emote FROM server WHERE server_id = ?")
	if err != nil {
		return "", err
	}

	parsed, err := strconv.ParseInt(guildId, 10, 0)
	if err != nil {
		return "", err
	}

	var archiveEmoji string
	if err = stmt.QueryRow(parsed).Scan(&archiveEmoji); err != nil {
		return "", err
	}

	return archiveEmoji, nil
}

func SetEmoji(guildId string, emote string) error {
	parsed, err := strconv.ParseInt(guildId, 10, 0)
	if err != nil {
		return err
	}

	if _, err := conn.Exec("UPDATE server SET archive_emote = ? WHERE server_id = ?", emote, parsed); err != nil {
		return err
	}

	return nil
}

func GetCount(guildId string) (int, error) {
	stmt, err := conn.Prepare("SELECT archive_emote_amount FROM server WHERE server_id = ?")
	if err != nil {
		return -1, err
	}

	parsed, err := strconv.ParseInt(guildId, 10, 0)
	if err != nil {
		return -1, err
	}

	var archiveCount int
	if err = stmt.QueryRow(parsed).Scan(&archiveCount); err != nil {
		return -1, err
	}

	return archiveCount, nil
}

func SetCount(guildId string, amount int64) error {
	parsed, err := strconv.ParseInt(guildId, 10, 0)
	if err != nil {
		return err
	}

	if _, err := conn.Exec("UPDATE server SET archive_emote_amount = ? WHERE server_id = ?", amount, parsed); err != nil {
		return err
	}

	return nil
}

func GetCustomCount(guildId string, channelId string) (int, error) {
	stmt, err := conn.Prepare("SELECT amount FROM custom_count WHERE server_id = ? AND channel_id = ?")
	if err != nil {
		return -1, err
	}

	parsedGuildId, err := strconv.ParseInt(guildId, 10, 0)
	if err != nil {
		return -1, err
	}
	parsedChannelId, err := strconv.ParseInt(channelId, 10, 0)
	if err != nil {
		return -1, err
	}

	var archiveCount int
	err = stmt.QueryRow(parsedGuildId, parsedChannelId).Scan(&archiveCount)
	if err == sql.ErrNoRows {
		return -1, nil
	} else if err != nil {
		return -1, err
	}

	return archiveCount, nil
}

func SetCustomCount(guildId string, channelId string, amount int64) error {
	count, err := GetCustomCount(guildId, channelId)
	if err != nil {
		return err
	}

	parsedGuildId, err := strconv.ParseInt(guildId, 10, 0)
	if err != nil {
		return err
	}
	parsedChannelId, err := strconv.ParseInt(channelId, 10, 0)
	if err != nil {
		return err
	}

	if count == -1 {
		if _, err := conn.Exec("INSERT INTO custom_count (server_id, channel_id, amount) VALUES (?, ?, ?)", parsedGuildId, parsedChannelId, amount); err != nil {
			return err
		}
	} else {
		if _, err := conn.Exec("UPDATE custom_count SET amount = ? WHERE server_id = ? AND channel_id = ?", amount, parsedGuildId, parsedChannelId); err != nil {
			return err
		}
	}

	return nil
}

func IsArchived(guildId string, channelId string, msgId string) (bool, error) {
	stmt, err := conn.Prepare("SELECT id FROM ignore_list WHERE server_id = ? AND channel_id = ? AND message_id = ?")
	if err != nil {
		return false, err
	}

	parsedGuildId, err := strconv.ParseInt(guildId, 10, 0)
	if err != nil {
		return false, err
	}
	parsedChannelId, err := strconv.ParseInt(channelId, 10, 0)
	if err != nil {
		return false, err
	}
	parsedMsgId, err := strconv.ParseInt(msgId, 10, 0)
	if err != nil {
		return false, err
	}

	var id int
	err = stmt.QueryRow(parsedGuildId, parsedChannelId, parsedMsgId).Scan(&id)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}

func Archive(guildId string, channelId string, msgId string) error {
	parsedGuildId, err := strconv.ParseInt(guildId, 10, 0)
	if err != nil {
		return err
	}
	parsedChannelId, err := strconv.ParseInt(channelId, 10, 0)
	if err != nil {
		return err
	}
	parsedMsgId, err := strconv.ParseInt(msgId, 10, 0)
	if err != nil {
		return err
	}

	if _, err := conn.Exec("INSERT INTO ignore_list (server_id, channel_id, message_id) VALUES (?, ?, ?)", parsedGuildId, parsedChannelId, parsedMsgId); err != nil {
		return err
	}

	return nil
}
