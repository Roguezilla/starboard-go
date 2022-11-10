package starboard

import (
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"starboard/cogs/galleries/instagram"
	"starboard/cogs/galleries/reddit"
	"starboard/sqldb"
	"starboard/utils"

	"github.com/anaskhan96/soup"
	"github.com/bwmarrin/discordgo"
)

type embedInfo struct {
	Flag         string
	Content      string
	MediaURL     string
	CustomAuthor *discordgo.User
}

var overrides = map[string]string{}

// this will never throw an error
var urlRegex, _ = regexp.Compile(`((?:https?):(?://)+(?:[\w\d_.~\-!*'();:@&=+$,/?#[\]]*))`)
var twitterRegex, _ = regexp.Compile(`https://(?:mobile.)?(vx)?twitter\.com/.+/status/\d+(?:/photo/(\d+))?`)

func buildEmbedInfo(s *discordgo.Session, m *discordgo.Message) embedInfo {
	e := embedInfo{
		Flag:    "message",
		Content: m.Content,
	}

	if URL, ok := overrides[m.GuildID+m.ChannelID+m.ID]; ok {
		isVideo := regexp.MustCompile(`.mp4|.mov|.webm`).MatchString(URL)

		if isVideo {
			e.Flag = "video"
		} else {
			e.Flag = "image"
		}
		e.Content = m.Content
		e.MediaURL = URL

		delete(overrides, m.GuildID+m.ChannelID+m.ID)
	} else {
		match := urlRegex.FindString(m.Content)
		if match != "" && len(m.Embeds) > 0 && len(m.Attachments) == 0 {
			parsedURL, _ := url.Parse(match)

			if strings.Contains(match, "deviantart.com") || strings.Contains(match, "tumblr.com") {
				if resp, err := soup.Get(match); err == nil {
					e.Flag = "image"
					e.Content = "[Source](" + match + ")\n" + strings.TrimSpace(strings.ReplaceAll(m.Content, match, ""))
					e.MediaURL = soup.HTMLParse(resp).Find("meta", "property", "og:image").Attrs()["content"]
				}
			} else if strings.Contains(match, "twitter.com") {
				urlData := twitterRegex.FindStringSubmatch(match)

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

				if parsedURL.Query().Get("u") != "" {
					if user, err := s.User(parsedURL.Query().Get("u")); err == nil {
						e.CustomAuthor = user
					}
				}
			} else if strings.Contains(match, "youtube.com") || strings.Contains(match, "youtu.be") {
				var videoID string
				if parsedURL.Query().Get("v") != "" {
					videoID = parsedURL.Query().Get("v")
				} else {
					videoID = strings.Split(match, "/")[2]
				}

				e.Flag = "image"
				e.Content = "[Source](" + match + ")\n" + strings.TrimSpace(strings.ReplaceAll(m.Content, match, ""))
				e.MediaURL = "https://img.youtube.com/vi/" + videoID + "/0.jpg"
			} else if strings.Contains(match, "imgur") {
				e.Flag = "image"
				if !strings.Contains(match, "i.imgur") {
					e.Content = "[Source](" + match + ")\n" + strings.TrimSpace(strings.ReplaceAll(m.Content, match, ""))
					if resp, err := soup.Get(match); err == nil {
						e.MediaURL = strings.ReplaceAll(soup.HTMLParse(resp).Find("meta", "property", "og:image").Attrs()["content"], "?fb", "")
					}
				} else {
					e.Content = strings.TrimSpace(strings.ReplaceAll(m.Content, match, ""))
					e.MediaURL = match
				}
			} else if regexp.MustCompile(`.mp4|.mov|.webm`).MatchString(match) {
				e.Flag = "video"
				e.Content = strings.TrimSpace(strings.ReplaceAll(m.Content, match, ""))
				e.MediaURL = match
			} else {
				if m.Embeds[0].Thumbnail != nil || m.Embeds[0].Image != nil {
					e.Flag = "image"
					e.Content = "[Source](" + match + ")\n" + strings.TrimSpace(strings.ReplaceAll(m.Content, match, ""))
					if m.Embeds[0].Thumbnail != nil {
						e.MediaURL = m.Embeds[0].Thumbnail.URL
					} else {
						e.MediaURL = m.Embeds[0].Image.URL
					}

					if parsedURL.Query().Get("u") != "" {
						if user, err := s.User(parsedURL.Query().Get("u")); err == nil {
							e.CustomAuthor = user
						}
					}
				}
			}
		} else {
			if len(m.Attachments) > 0 {
				isVideo := regexp.MustCompile(`.mp4|.mov|.webm`).MatchString(m.Attachments[0].URL)
				isSpoiler := strings.HasPrefix(m.Attachments[0].Filename, "SPOILER_")

				if isVideo {
					e.Flag = "video"
				} else {
					e.Flag = "image"
				}

				e.Content = m.Content
				if isSpoiler {
					e.MediaURL = "https://i.imgur.com/GFn7HTJ.png"
				} else {
					e.MediaURL = m.Attachments[0].URL
				}
			} else {
				if reddit.ValidateEmbed(m.Embeds) || instagram.ValidateEmbed(m.Embeds) {
					e.Flag = "image"
					e.MediaURL = m.Embeds[0].Image.URL
					if user, err := s.User(m.Embeds[0].Fields[0].Value[2 : len(m.Embeds[0].Fields[0].Value)-1]); err == nil {
						e.CustomAuthor = user
					}
				}
			}
		}
	}

	return e
}

func HandleMessageReactionAdd(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
	archived, err := sqldb.IsArchived(m.GuildID, m.ChannelID, m.MessageID)
	if err != nil {
		s.ChannelMessageSendReply(m.ChannelID, err.Error(), &discordgo.MessageReference{GuildID: m.GuildID, ChannelID: m.ChannelID, MessageID: m.MessageID})
		return
	}
	if archived {
		return
	}

	emoji, err := sqldb.Emoji(m.GuildID)
	if err != nil {
		s.ChannelMessageSendReply(m.ChannelID, err.Error(), &discordgo.MessageReference{GuildID: m.GuildID, ChannelID: m.ChannelID, MessageID: m.MessageID})
		return
	} else if emoji != m.Emoji.APIName() {
		return
	}

	msg, err := s.ChannelMessage(m.ChannelID, m.MessageID)
	if err != nil {
		s.ChannelMessageSendReply(m.ChannelID, err.Error(), &discordgo.MessageReference{GuildID: m.GuildID, ChannelID: m.ChannelID, MessageID: m.MessageID})
		return
	}
	msg.GuildID = m.GuildID // ChannelMessage returns Message with empty GuildID

	amount, err := sqldb.ChannelAmount(m.GuildID, m.ChannelID)
	if err != nil {
		s.ChannelMessageSendReply(m.ChannelID, err.Error(), &discordgo.MessageReference{GuildID: m.GuildID, ChannelID: m.ChannelID, MessageID: m.MessageID})
		return
	}

	if amount == -1 {
		amount, err = sqldb.GlobalAmount(m.GuildID)
		if err != nil {
			s.ChannelMessageSendReply(m.ChannelID, err.Error(), &discordgo.MessageReference{GuildID: m.GuildID, ChannelID: m.ChannelID, MessageID: m.MessageID})
			return
		}
	}

	emojiCount, err := utils.EmojiCount(s, m)
	if err != nil {
		s.ChannelMessageSendReply(m.ChannelID, err.Error(), &discordgo.MessageReference{GuildID: m.GuildID, ChannelID: m.ChannelID, MessageID: m.MessageID})
		return
	}

	if emojiCount >= amount {
		channelID, err := sqldb.Channel(m.GuildID)
		if err != nil {
			s.ChannelMessageSendReply(m.ChannelID, err.Error(), &discordgo.MessageReference{GuildID: m.GuildID, ChannelID: m.ChannelID, MessageID: m.MessageID})
			return
		}

		embedInfo := buildEmbedInfo(s, msg)

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
				Name:    msg.Author.Username,
				IconURL: msg.Author.AvatarURL(""),
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
			Value:  "[Jump to ](https://discordapp.com/channels/" + msg.GuildID + "/" + msg.ChannelID + "/" + msg.ID + ")<#" + msg.ChannelID + ">",
			Inline: false,
		})

		if err := sqldb.Archive(msg.GuildID, msg.ChannelID, msg.ID); err != nil {
			s.ChannelMessageSendReply(msg.ChannelID, err.Error(), msg.Reference())
			return
		}

		if _, err := s.ChannelMessageSendEmbed(channelID, &embed); err != nil {
			s.ChannelMessageSendReply(msg.ChannelID, err.Error(), msg.Reference())
			return
		}

		if embedInfo.Flag == "video" && embedInfo.MediaURL != "" {
			if _, err := s.ChannelMessageSend(channelID, embedInfo.MediaURL); err != nil {
				s.ChannelMessageSendReply(msg.ChannelID, err.Error(), msg.Reference())
				return
			}
		}
	}
}

func ArchiveOverrideCommand(s *discordgo.Session, m *discordgo.MessageCreate, numArgs int, args ...string) {
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
		s.ChannelMessageSendReply(m.ChannelID, "You have to reply to the message you want to override.", m.Message.Reference())
		return
	}

	if _, ok := overrides[m.GuildID+m.ReferencedMessage.ChannelID+m.ReferencedMessage.ID]; !ok {
		overrides[m.GuildID+m.ReferencedMessage.ChannelID+m.ReferencedMessage.ID] = args[0]
	} else {
		s.ChannelMessageSendReply(m.ChannelID, "Message already has override.", m.Message.Reference())
	}

	s.ChannelMessageDelete(m.ChannelID, m.ID)
}
