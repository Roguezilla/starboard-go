package sqldb

import (
	"database/sql"

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

func Token() (string, error) {
	stmt, err := conn.Prepare("SELECT value FROM settings WHERE key = 'token'")
	if err != nil {
		return "", err
	}

	var token string
	if err = stmt.QueryRow().Scan(&token); err != nil {
		return "", err
	}

	return token, nil
}

func Setup(guildID string, channelID string, emoji string, amount int64) error {
	if _, err := conn.Exec("INSERT INTO guild_data (guild, channel, emoji, amount) VALUES (?, ?, ?, ?)", guildID, channelID, emoji, amount); err != nil {
		return err
	}

	return nil
}

func IsSetup(guildID string) (bool, error) {
	_, err := Emoji(guildID)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}

func Channel(guildID string) (string, error) {
	stmt, err := conn.Prepare("SELECT channel FROM guild_data WHERE guild = ?")
	if err != nil {
		return "", err
	}

	var archiveChannel string
	if err = stmt.QueryRow(guildID).Scan(&archiveChannel); err != nil {
		return "", err
	}

	return archiveChannel, nil
}

func SetChannel(guildID string, channelID string) error {
	if _, err := conn.Exec("UPDATE guild_data SET channel = ? WHERE guild = ?", channelID, guildID); err != nil {
		return err
	}

	return nil
}

func Emoji(guildID string) (string, error) {
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

func SetEmoji(guildID string, emoji string) error {
	if _, err := conn.Exec("UPDATE guild_data SET emoji = ? WHERE guild = ?", emoji, guildID); err != nil {
		return err
	}

	return nil
}

func GlobalAmount(guildID string) (int, error) {
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

func SetAmount(guildID string, amount int64) error {
	if _, err := conn.Exec("UPDATE guild_data SET amount = ? WHERE guild = ?", amount, guildID); err != nil {
		return err
	}

	return nil
}

func ChannelAmount(guildID string, channelID string) (int, error) {
	stmt, err := conn.Prepare("SELECT amount FROM custom_channel_amount WHERE guild = ? AND channel = ?")
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

func SetChannelAmount(guildID string, channelID string, amount int64) error {
	channelAmount, err := ChannelAmount(guildID, channelID)
	if err != nil {
		return err
	}

	if channelAmount == -1 {
		if _, err := conn.Exec("INSERT INTO custom_channel_amount (guild, channel, amount) VALUES (?, ?, ?)", guildID, channelID, amount); err != nil {
			return err
		}
	} else {
		if _, err := conn.Exec("UPDATE custom_channel_amount SET amount = ? WHERE guild = ? AND channel = ?", amount, guildID, channelID); err != nil {
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

func Unarchive(guildID string, channelID string, msgID string) error {
	if _, err := conn.Exec("DELETE FROM archive WHERE guild = ? AND channel = ? AND message = ?", guildID, channelID, msgID); err != nil {
		return err
	}

	return nil
}
