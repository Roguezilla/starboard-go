package reddit

import (
	"crypto/tls"
	"encoding/json"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type gallery struct {
	CurrentIdx int
	Gallery    []string
}

var redditCache = map[string]*gallery{}

var regex, _ = regexp.Compile(`^((?:(?:(?:https):(?://)+)(?:www\.)?)redd(?:it\.com/|\.it/).+)$`)

func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	match := regex.FindStringSubmatch(m.Content)
	if len(match) == 0 {
		return
	}

	imageLink, title, err := postData(match[0], m.Message)
	if err != nil {
		s.ChannelMessageSendReply(m.ChannelID, err.Error(), &discordgo.MessageReference{GuildID: m.GuildID, ChannelID: m.ChannelID, MessageID: m.ID})
		return
	}

	embed := discordgo.MessageEmbed{
		Title: title,
		URL:   match[1],
		Color: 0xffcc00,
		Image: &discordgo.MessageEmbedImage{
			URL: imageLink,
		},
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Original Poster",
				Value:  m.Author.Mention(),
				Inline: true,
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "by rogue#0001",
		},
	}

	msg, err := s.ChannelMessageSendEmbed(m.ChannelID, &embed)
	if err != nil {
		return
	}

	if entry, ok := redditCache[m.GuildID+m.ChannelID+m.ID]; ok {
		// transfer the gallery cache of the original message to the new one(the embed)
		redditCache[msg.GuildID+msg.ChannelID+msg.ID] = &gallery{
			CurrentIdx: entry.CurrentIdx,
			Gallery:    entry.Gallery,
		}

		msg.Embeds[0].Fields = append(msg.Embeds[0].Fields, &discordgo.MessageEmbedField{
			Name:   "Page",
			Value:  strconv.Itoa(entry.CurrentIdx+1) + "/" + strconv.Itoa(len(entry.Gallery)),
			Inline: true,
		})

		s.ChannelMessageEditEmbeds(msg.ChannelID, msg.ID, msg.Embeds)

		s.MessageReactionAdd(msg.ChannelID, msg.ID, string([]rune{11013, 65039}))
		s.MessageReactionAdd(msg.ChannelID, msg.ID, string([]rune{10145, 65039}))

		// delete original message, as we don't need it anymore
		delete(redditCache, m.GuildID+m.ChannelID+m.ID)
	}

	s.ChannelMessageDelete(m.ChannelID, m.ID)
}

func MessageReactionAdd(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
	// right and left arrow emojis in rune format
	if m.Emoji.APIName() != string([]rune{10145, 65039}) && m.Emoji.APIName() != string([]rune{11013, 65039}) {
		return
	}

	msg, err := s.ChannelMessage(m.ChannelID, m.MessageID)
	if err != nil {
		s.ChannelMessageSendReply(m.ChannelID, err.Error(), &discordgo.MessageReference{GuildID: m.GuildID, ChannelID: m.ChannelID, MessageID: m.MessageID})
		return
	}
	msg.GuildID = m.GuildID // ChannelMessage returns Message with empty GuildID

	if !ValidateEmbed(msg.Embeds) {
		return
	}

	// the gallery cache gets wiped when the bot is turned off, so we have to rebuild it
	if _, ok := redditCache[msg.GuildID+msg.ChannelID+msg.ID]; !ok {
		api, err := urlData(msg.Embeds[0].URL)
		if err != nil {
			s.ChannelMessageSendReply(m.ChannelID, err.Error(), &discordgo.MessageReference{GuildID: m.GuildID, ChannelID: m.ChannelID, MessageID: m.MessageID})
			return
		} else if !buildGallery(api, msg, true) {
			return
		}
	}

	if entry, ok := redditCache[msg.GuildID+msg.ChannelID+msg.ID]; ok {
		// right arrow
		if m.Emoji.APIName() == string([]rune{10145, 65039}) {
			entry.CurrentIdx++
			if entry.CurrentIdx > len(entry.Gallery)-1 {
				entry.CurrentIdx = 0
			}
		} else {
			entry.CurrentIdx--
			if entry.CurrentIdx < 0 {
				entry.CurrentIdx = len(entry.Gallery) - 1
			}
		}

		msg.Embeds[0].Image.URL = entry.Gallery[entry.CurrentIdx]
		msg.Embeds[0].Fields[1].Value = strconv.Itoa(entry.CurrentIdx+1) + "/" + strconv.Itoa(len(entry.Gallery))
		s.ChannelMessageEditEmbeds(msg.ChannelID, msg.ID, msg.Embeds)

		s.MessageReactionRemove(m.ChannelID, m.MessageID, m.Emoji.APIName(), m.UserID)
	}
}

func buildGallery(api AutoGenerated, m *discordgo.Message, rebuild bool) bool {
	redditCache[m.GuildID+m.ChannelID+m.ID] = &gallery{}

	if entry, ok := redditCache[m.GuildID+m.ChannelID+m.ID]; ok {
		for _, item := range api[0].Data.Children[0].Data.GalleryData.Items {
			entry.Gallery = append(entry.Gallery, strings.ReplaceAll(mediaMetadataByID(api[0].Data.Children[0].Data.MediaMetadata, item.MediaID), "&amp;", "&"))
		}
	}

	// when rebuilding, match current index in the rebuilt cache to the current image in the embed
	if rebuild {
		if entry, ok := redditCache[m.GuildID+m.ChannelID+m.ID]; ok {
			for i := 0; i < len(entry.Gallery); i++ {
				if entry.Gallery[i] == m.Embeds[0].Image.URL {
					entry.CurrentIdx = i
					break
				}
			}
		}
	}

	return true
}

func ValidateEmbed(embeds []*discordgo.MessageEmbed) bool {
	if len(embeds) == 0 {
		return false
	}

	match := regex.FindStringSubmatch(embeds[0].URL)
	return len(match) != 0 && len(embeds[0].Fields) == 2 && embeds[0].Fields[1].Name == "Page"
}

func urlData(url string) (AutoGenerated, error) {
	// get the final url from a redd.it redirect
	if strings.Contains(url, "redd.it") {
		res, err := http.Head(url)
		if err != nil {
			return nil, err
		}

		url = res.Request.URL.String() + ".json"
	} else {
		url = strings.Split(url, "?")[0] + ".json"
	}

	// https://www.reddit.com/r/redditdev/comments/uncu00/comment/i8gyfm
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{},
		},
	}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-agent", "starboard-go")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	resData, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	apiJSON := AutoGenerated{}
	err = json.Unmarshal(resData, &apiJSON)
	if err != nil {
		return nil, err
	}

	return apiJSON, nil
}

func postData(url string, msg *discordgo.Message) (string, string, error) {
	api, err := urlData(url)
	if err != nil {
		return "", "", err
	}

	if api[0].Data.Children[0].Data.IsGallery {
		// iirc MediaMetadata can be unordered and it's plain easier to access the images with the help of GalleryData(which is ordered for sure)
		url = mediaMetadataByID(api[0].Data.Children[0].Data.MediaMetadata, api[0].Data.Children[0].Data.GalleryData.Items[0].MediaID)

		buildGallery(api, msg, false)
	} else {
		// URLOverriddenByDest for videos is the link to the post, so we have to use a preview image
		if api[0].Data.Children[0].Data.IsVideo {
			url = api[0].Data.Children[0].Data.Preview.Images[0].Source.URL
		} else {
			// covers pretty much everything that's not a video
			url = api[0].Data.Children[0].Data.URLOverriddenByDest
		}
	}

	return strings.ReplaceAll(url, "&amp;", "&"), api[0].Data.Children[0].Data.Title, nil
}
