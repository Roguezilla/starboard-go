package starboard

import (
	"fmt"
	"log"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"starboard/cogs/galleries/instagram"
	"starboard/cogs/galleries/reddit"
	"starboard/sqldb"
	"starboard/utils"

	"github.com/bwmarrin/discordgo"
)

type embedInfo struct {
	Flag     string
	Content  string
	MediaURL string
	Author   *discordgo.User
}

var overrides = map[string]string{}

// this will never throw an error
var urlRegex, _ = regexp.Compile(`https?://[\w\d_.~\-!*'();:@&=+$,/?#[\]]*`)
var twitterRegex, _ = regexp.Compile(`https?://vxtwitter\.com/.+/status/\d+`)
var youtubeRegex, _ = regexp.Compile(`https?://(?:(?:www\.)?youtube.com/watch\?v=|youtu\.be/)[A-Za-z0-9_\-]{11}`)
var imgurRegex, _ = regexp.Compile(`https?://(?:i\.)?imgur.com/(?:gallery/.+|.+\..+)`)
var tenorRegex, _ = regexp.Compile(`https?://tenor\.com/view/.+`)

func buildEmbedInfo(s *discordgo.Session, m *discordgo.Message) embedInfo {
	e := embedInfo{
		Flag:    "message",
		Content: m.Content,
		Author:  m.Author,
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

			if twitterRegex.FindString(match) != "" {
				if m.Embeds[0].Video != nil {
					e.Flag = "video"
					e.MediaURL = m.Embeds[0].Video.URL
				} else {
					e.Flag = "image"
					e.MediaURL = m.Embeds[0].Thumbnail.URL
				}

				e.Content = "[" + m.Embeds[0].Title + "](" + match + ")\n\n" + m.Embeds[0].Description

				if parsedURL.Query().Get("u") != "" {
					if user, err := s.User(parsedURL.Query().Get("u")); err == nil {
						e.Author = user
					}
				}
			} else if youtubeRegex.FindString(match) != "" {
				videoID := ""
				if parsedURL.Query().Get("v") != "" {
					videoID = parsedURL.Query().Get("v")
				} else {
					videoID = strings.Split(match, "youtu.be/")[1][0:11]
					fmt.Printf("videoID: %v\n", videoID)
				}

				e.Flag = "image"

				e.Content = "[Source](" + match + ")\n" + strings.TrimSpace(strings.ReplaceAll(m.Content, match, ""))

				e.MediaURL = "https://img.youtube.com/vi/" + videoID + "/0.jpg"
			} else if regexp.MustCompile(`.mp4|.mov|.webm`).MatchString(match) {
				e.Flag = "video"

				e.Content = "[The video below](https://youtu.be/dQw4w9WgXcQ)\n" + strings.TrimSpace(strings.ReplaceAll(m.Content, match, ""))
				fmt.Printf("e.Content: %v\n", e.Content)

				e.MediaURL = match
			} else if m.Embeds[0].Thumbnail != nil || m.Embeds[0].Image != nil {
				e.Flag = "image"

				e.Content = "[Source](" + match + ")\n" + strings.TrimSpace(strings.ReplaceAll(m.Content, match, ""))

				if m.Embeds[0].Thumbnail != nil {
					e.MediaURL = m.Embeds[0].Thumbnail.URL
					if imgurRegex.FindString(match) != "" {
						e.MediaURL = m.Embeds[0].Thumbnail.ProxyURL
					} else if tenorRegex.FindString(match) != "" {
						splitUrl := strings.Split(m.Embeds[0].Thumbnail.URL, "")
						splitUrl[39] = strings.ToLower(splitUrl[39])
						e.MediaURL = strings.ReplaceAll(strings.Join(splitUrl, ""), ".png", ".gif")
					}
				} else {
					e.MediaURL = m.Embeds[0].Image.URL
				}

				if parsedURL.Query().Get("u") != "" {
					if user, err := s.User(parsedURL.Query().Get("u")); err == nil {
						e.Author = user
					}
				}
			}
		} else {
			if len(m.Attachments) > 0 {
				isVideo := regexp.MustCompile(`.mp4|.mov|.webm`).MatchString(m.Attachments[0].URL)
				isSpoiler := strings.HasPrefix(m.Attachments[0].Filename, "SPOILER_")

				if isVideo {
					e.Flag = "video"
					if isSpoiler {
						e.Flag = "image"
					}
				} else {
					e.Flag = "image"
				}

				e.Content = m.Content
				if isVideo {
					e.Content = "[The video below](https://youtu.be/dQw4w9WgXcQ)\n" + strings.TrimSpace(strings.ReplaceAll(m.Content, match, ""))
				}

				e.MediaURL = m.Attachments[0].URL
				if isSpoiler {
					e.MediaURL = "https://i.imgur.com/GFn7HTJ.png"
				}
			} else {
				if reddit.ValidateEmbed(m.Embeds) || instagram.ValidateEmbed(m.Embeds) {
					e.Flag = "image"
					e.MediaURL = m.Embeds[0].Image.URL
					if user, err := s.User(m.Embeds[0].Fields[0].Value[2 : len(m.Embeds[0].Fields[0].Value)-1]); err == nil {
						e.Author = user
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
		log.Println("starboard.go HandleMessageReactionAdd sqldb.IsArchived:", err)
		return
	}
	if archived {
		return
	}

	emoji, err := sqldb.Emoji(m.GuildID)
	if err != nil {
		log.Println("starboard.go HandleMessageReactionAdd sqldb.Emoji:", err)
		return
	} else if emoji != m.Emoji.APIName() {
		return
	}

	msg, err := s.ChannelMessage(m.ChannelID, m.MessageID)
	if err != nil {
		log.Println("starboard.go HandleMessageReactionAdd s.ChannelMessage:", err)
		return
	}
	msg.GuildID = m.GuildID // ChannelMessage returns Message with empty GuildID

	amount, err := sqldb.ChannelAmount(m.GuildID, m.ChannelID)
	if err != nil {
		log.Println("starboard.go HandleMessageReactionAdd sqldb.ChannelAmount:", err)
		return
	}

	if amount == -1 {
		amount, err = sqldb.GlobalAmount(m.GuildID)
		if err != nil {
			log.Println("starboard.go HandleMessageReactionAdd sqldb.GlobalAmount:", err)
			return
		}
	}

	emojiCount, err := utils.EmojiCount(s, m)
	if err != nil {
		log.Println("starboard.go HandleMessageReactionAdd utils.EmojiCount:", err)
		return
	}

	if emojiCount >= amount {
		channelID, err := sqldb.Channel(m.GuildID)
		if err != nil {
			log.Println("starboard.go HandleMessageReactionAdd sqldb.Channel:", err)
			return
		}

		embedInfo := buildEmbedInfo(s, msg)

		embed := discordgo.MessageEmbed{
			Color: 0xffcc00,
			Footer: &discordgo.MessageEmbedFooter{
				Text: "by rogue#0001",
			},
		}

		embed.Author = &discordgo.MessageEmbedAuthor{
			Name:    embedInfo.Author.Username,
			IconURL: embedInfo.Author.AvatarURL(""),
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
			log.Println("starboard.go HandleMessageReactionAdd sqldb.Archive:", err)
			return
		}

		if _, err := s.ChannelMessageSendEmbed(channelID, &embed); err != nil {
			log.Println("starboard.go HandleMessageReactionAdd s.ChannelMessageSendEmbed:", err)
			return
		}

		if embedInfo.Flag == "video" && embedInfo.MediaURL != "" {
			if _, err := s.ChannelMessageSend(channelID, embedInfo.MediaURL); err != nil {
				log.Println("starboard.go HandleMessageReactionAdd s.ChannelMessageSend:", err)
				return
			}
		}
	}
}

func ArchiveOverrideCommand(s *discordgo.Session, m *discordgo.MessageCreate, numArgs int, args ...string) {
	if numArgs != len(args) {
		s.ChannelMessageSend(m.ChannelID, "❌Invalid number of arguments, got "+strconv.Itoa(len(args))+" expected "+strconv.Itoa(numArgs)+".")
		return
	}

	if has, err := utils.CheckPermission(s, m.Message, discordgo.PermissionManageMessages); !has {
		if err != nil {
			log.Println("starboard.go ArchiveOverrideCommand tils.CheckPermission:", err)
			return
		} else {
			s.ChannelMessageSendReply(m.ChannelID, "❌You don't have permission to do that.", m.Message.Reference())
			return
		}
	}
	/* */
	if m.MessageReference == nil {
		s.ChannelMessageSendReply(m.ChannelID, "❌You have to reply to the message you want to override.", m.Message.Reference())
		return
	}

	if _, ok := overrides[m.GuildID+m.ReferencedMessage.ChannelID+m.ReferencedMessage.ID]; !ok {
		overrides[m.GuildID+m.ReferencedMessage.ChannelID+m.ReferencedMessage.ID] = args[0]
	} else {
		s.ChannelMessageSendReply(m.ChannelID, "❌Message already has override.", m.Message.Reference())
	}

	s.ChannelMessageDelete(m.ChannelID, m.ID)
}
