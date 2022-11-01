package cogs

import (
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/anaskhan96/soup"
	"github.com/bwmarrin/discordgo"
	"roguezilla.github.io/starboard/sqldb"
)

type embedInfo struct {
	Flag         string
	Content      string
	MediaURL     string
	CustomAuthor *discordgo.User
}

// this will never throw an error
var urlRegex, _ = regexp.Compile(`((?:https?):(?://)+(?:[\w\d_.~\-!*'();:@&=+$,/?#[\]]*))`)
var twitterRegex, _ = regexp.Compile(`https://(?:mobile.)?(vx)?twitter\.com/.+/status/\d+(?:/photo/(\d+))?`)

func buildEmbedInfo(s *discordgo.Session, m *discordgo.Message) embedInfo {
	e := embedInfo{
		Flag:    "message",
		Content: m.Content,
	}

	match := urlRegex.FindString(m.Content)
	if match != "" && len(m.Embeds) > 0 && len(m.Attachments) == 0 {
		if strings.Contains(match, "deviantart.com") || strings.Contains(match, "tumblr.com") {
			resp, err := soup.Get(match)
			if err != nil {
				return e
			}

			e.Flag = "image"
			e.Content = "[Source](" + match + ")\n" + strings.TrimSpace(strings.ReplaceAll(m.Content, match, ""))
			e.MediaURL = soup.HTMLParse(resp).Find("meta", "property", "og:image").Attrs()["content"]
		} else if strings.Contains(match, "twitter.com") {
			urlData := twitterRegex.FindStringSubmatch(match)
			u, _ := url.Parse(match)

			if urlData[1] != "" && m.Embeds[0].Video != nil {
				e.Flag = "video"
				e.MediaURL = m.Embeds[0].Video.URL
			} else if m.Embeds[0].Image != nil {
				e.Flag = "image"
				e.MediaURL = m.Embeds[0].Image.URL
				if urlData[2] != "" {
					idx, err := strconv.ParseInt(urlData[2], 10, 64)
					if err == nil {
						e.MediaURL = m.Embeds[idx-1].Image.URL
					}
				}
			}

			e.Content = "[Tweet](" + match + ")\n" + strings.TrimSpace(strings.ReplaceAll(m.Content, match, ""))

			if u.Query().Get("u") != "" {
				if user, err := s.User(u.Query().Get("u")); err == nil {
					e.CustomAuthor = user
				}
			}
		}
	} else {

	}

	return e
}

func Archive(s *discordgo.Session, m *discordgo.Message, channelID string) {
	embedInfo := buildEmbedInfo(s, m)

	embed := discordgo.MessageEmbed{
		Color: 0xffcc00,
		Footer: &discordgo.MessageEmbedFooter{
			Text: "by rogue#0001",
		},
	}

	if embedInfo.CustomAuthor != nil {
		embed.Author = &discordgo.MessageEmbedAuthor{
			Name:    embedInfo.CustomAuthor.Username,
			IconURL: embedInfo.CustomAuthor.AvatarURL(""),
		}
	} else {
		embed.Author = &discordgo.MessageEmbedAuthor{
			Name:    m.Author.Username,
			IconURL: m.Author.AvatarURL(""),
		}
	}
	if embedInfo.Content != "" {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "What?",
			Value:  embedInfo.Content,
			Inline: false,
		})
	}
	if embedInfo.Flag == "image" && embedInfo.MediaURL != "" {
		embed.Image = &discordgo.MessageEmbedImage{
			URL: embedInfo.MediaURL,
		}
	}

	embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
		Name:   "Where?",
		Value:  "[Jump to ](https://discordapp.com/channels/" + m.GuildID + "/" + m.ChannelID + "/" + m.ID + ")<#" + m.ChannelID + ">",
		Inline: false,
	})

	if err := sqldb.Archive(m.GuildID, m.ChannelID, m.ID); err != nil {
		s.ChannelMessageSendReply(m.ChannelID, err.Error(), m.Reference())
		return
	}

	if _, err := s.ChannelMessageSendEmbed(channelID, &embed); err != nil {
		s.ChannelMessageSendReply(m.ChannelID, err.Error(), m.Reference())
		return
	}

	if embedInfo.Flag == "video" && embedInfo.MediaURL != "" {
		if _, err := s.ChannelMessageSend(channelID, embedInfo.MediaURL); err != nil {
			s.ChannelMessageSendReply(m.ChannelID, err.Error(), m.Reference())
			return
		}
	}
}
