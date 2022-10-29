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
	rows, err := conn.Query("SELECT value FROM settings WHERE key = 'token'")
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

func IsSetup(guildID string) (bool, error) {
	_, err := GetEmoji(guildID)
	if err != nil {
		return false, err
	}

	return true, nil
}

func Setup(guildID string, channelID string, emoji string, amount int64) (bool, error) {
	if _, err := conn.Exec("INSERT INTO guild_data (guild, channel, emoji, amount) VALUES (?, ?, ?, ?)", guildID, channelID, emoji, amount); err != nil {
		return false, err
	}

	return true, nil
}

func GetChannel(guildID string) (int, error) {
	stmt, err := conn.Prepare("SELECT channel FROM guild_data WHERE guild = ?")
	if err != nil {
		return -1, err
	}

	parsed, err := strconv.ParseInt(guildID, 10, 0)
	if err != nil {
		return -1, err
	}

	var archiveChannel int
	if err = stmt.QueryRow(parsed).Scan(&archiveChannel); err != nil {
		return -1, err
	}

	return archiveChannel, nil
}

func SetChannel(guildID string, channel int) error {
	parsed, err := strconv.ParseInt(guildID, 10, 0)
	if err != nil {
		return err
	}

	if _, err := conn.Exec("UPDATE guild_data SET channel = ? WHERE guild = ?", channel, parsed); err != nil {
		return err
	}

	return nil
}

func GetEmoji(guildID string) (string, error) {
	stmt, err := conn.Prepare("SELECT emoji FROM guild_data WHERE guild = ?")
	if err != nil {
		return "", err
	}

	var archiveEmoji string
	if err = stmt.QueryRow(guildID).Scan(&archiveEmoji); err != nil {
		return "", err
	}

	return archiveEmoji, nil
}

func SetEmoji(guildID string, emote string) error {
	parsed, err := strconv.ParseInt(guildID, 10, 0)
	if err != nil {
		return err
	}

	if _, err := conn.Exec("UPDATE guild_data SET emoji = ? WHERE guild = ?", emote, parsed); err != nil {
		return err
	}

	return nil
}

func GetCount(guildID string) (int, error) {
	stmt, err := conn.Prepare("SELECT amount FROM guild_data WHERE guild = ?")
	if err != nil {
		return -1, err
	}

	var amount int
	if err = stmt.QueryRow(guildID).Scan(&amount); err != nil {
		return -1, err
	}

	return amount, nil
}

func SetCount(guildID string, amount int64) error {
	if _, err := conn.Exec("UPDATE guild_data SET amount = ? WHERE guild = ?", amount, guildID); err != nil {
		return err
	}

	return nil
}

func GetCustomCount(guildID string, channelID string) (int, error) {
	stmt, err := conn.Prepare("SELECT amount FROM custom_count WHERE guild = ? AND channel = ?")
	if err != nil {
		return -1, err
	}

	var amount int
	err = stmt.QueryRow(guildID, channelID).Scan(&amount)
	if err == sql.ErrNoRows {
		return -1, nil
	} else if err != nil {
		return -1, err
	}

	return amount, nil
}

func SetCustomCount(guildID string, channelID string, amount int64) error {
	count, err := GetCustomCount(guildID, channelID)
	if err != nil {
		return err
	}

	if count == -1 {
		if _, err := conn.Exec("INSERT INTO custom_count (guild, channel, amount) VALUES (?, ?, ?)", guildID, channelID, amount); err != nil {
			return err
		}
	} else {
		if _, err := conn.Exec("UPDATE custom_count SET amount = ? WHERE guild = ? AND channel = ?", amount, guildID, channelID); err != nil {
			return err
		}
	}

	return nil
}

func IsArchived(guildID string, channelID string, msgID string) (bool, error) {
	stmt, err := conn.Prepare("SELECT id FROM archive WHERE guild = ? AND channel = ? AND message = ?")
	if err != nil {
		return false, err
	}

	var id int
	err = stmt.QueryRow(guildID, channelID, msgID).Scan(&id)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}

func Archive(guildID string, channelID string, msgID string) error {
	if _, err := conn.Exec("INSERT INTO archive (guild, channel, message) VALUES (?, ?, ?)", guildID, channelID, msgID); err != nil {
		return err
	}

	return nil
}
