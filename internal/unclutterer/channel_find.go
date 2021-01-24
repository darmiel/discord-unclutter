package unclutterer

import (
	"errors"
	"github.com/bwmarrin/discordgo"
	duconfig "github.com/darmiel/discord-unclutterer/internal/unclutterer/config"
	"log"
)

var channelCache = make(map[string]*ChannelCategory)

func FindCreatedTextChannels(s *discordgo.Session, guildID string, config *duconfig.Config) (channels []*UnclutteredChannel, category *discordgo.Channel, err error) {
	allChannels, err := s.GuildChannels(guildID)
	if err != nil {
		return nil, nil, err
	}

	for _, c := range allChannels {
		// check for category
		if c.Type == discordgo.ChannelTypeGuildCategory {
			if c.Name == config.ChannelCategoryName {
				category = c
				continue
			}
		}

		// check if text channel
		if c.Type == discordgo.ChannelTypeGuildText {
			if id := extractChannelID(c, config); id != "" {
				channels = append(channels, &UnclutteredChannel{
					Channel:        c,
					VoiceChannelID: id,
				})
			}
		}
	}

	return
}

func (us *UserVoiceStateSession) findOrCreateText(voiceChannelID string) (chcat *ChannelCategory, err error) {
	if voiceChannelID == "" {
		return nil, errors.New("empty voice channel id")
	}

	var chID string
	var catID string

	chcat, ok := channelCache[voiceChannelID]
	if ok {
		chID = chcat.Channel.ID
		catID = chcat.Category.ID
	} else {
		channels, category, _ := FindCreatedTextChannels(us.Session, us.GuildID, us.config)

		if category != nil {
			catID = category.ID
		}

		if channels != nil {
			for _, c := range channels {
				if c.VoiceChannelID == voiceChannelID {
					chID = c.ID
				}
			}
		}
	}

	var chChannel *discordgo.Channel
	var catChannel *discordgo.Channel

	// check if channels exist
	if chC, err := us.Session.State.Channel(chID); err != nil {
		if err != discordgo.ErrStateNotFound {
			return nil, err
		}
		log.Println("🧠 Advanced ChID-Lookup")
		if chC, err := us.Session.Channel(chID); err == nil {
			chChannel = chC
		}
	} else {
		chChannel = chC
	}
	if catC, err := us.Session.State.Channel(catID); err != nil {
		if err != discordgo.ErrStateNotFound {
			return nil, err
		}
		log.Println("🧠 Advanced CatID-Lookup")
		if catC, err := us.Session.Channel(catID); err == nil {
			catChannel = catC
		}
	} else {
		catChannel = catC
	}

	// create channels, if necessary
	if catChannel == nil {
		if channel, err := us.CreateCategory(); err != nil {
			log.Println("ERROR: (creating category)", err)
			return nil, err
		} else {
			catChannel = channel
		}
	}

	if chChannel == nil {
		if channel, err := us.CreateChannel(catChannel.ID); err != nil {
			log.Println("ERROR: (creating text channel)", err)
			return nil, err
		} else {
			chChannel = channel
		}
	}

	if ok && chcat != nil {
		if chcat.Channel.ID == chChannel.ID && chcat.Category.ID == catChannel.ID {
			return chcat, nil
		}
	}

	c := &ChannelCategory{
		Channel:  chChannel,
		Category: catChannel,
	}

	// cache channel
	channelCache[voiceChannelID] = c

	return c, nil
}
