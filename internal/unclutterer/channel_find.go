package unclutterer

import (
	"errors"
	"github.com/bwmarrin/discordgo"
	"log"
)

var channelCache = make(map[string]*ChannelCategory)

func FindCreatedTextChannels(s *discordgo.Session, guildID string) (channels []*UnclutteredChannel, category *discordgo.Channel, err error) {
	allChannels, err := s.GuildChannels(guildID)
	if err != nil {
		return nil, nil, err
	}

	for _, c := range allChannels {
		// check for category
		if c.Type == discordgo.ChannelTypeGuildCategory {
			if c.Name == CategoryName {
				category = c
				continue
			}
		}

		// check if text channel
		if c.Type == discordgo.ChannelTypeGuildText {
			if id := extractChannelID(c); id != "" {
				channels = append(channels, &UnclutteredChannel{
					Channel:        c,
					VoiceChannelID: id,
				})
			}
		}
	}
	return
}

func (us *UserSess) findOrCreateText(voiceChannelID string) (chcat *ChannelCategory, err error) {
	if voiceChannelID == "" {
		return nil, errors.New("empty voice channel id")
	}

	var chID string
	var catID string

	chcat, ok := channelCache[voiceChannelID]
	if ok {
		chID = chcat.Channel.ID
		catID = chcat.Channel.ID
	} else {
		channels, category, _ := FindCreatedTextChannels(us.Session, us.GuildID)

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
	} else {
		chChannel = chC
	}
	if catC, err := us.Session.State.Channel(catID); err != nil {
		if err != discordgo.ErrStateNotFound {
			return nil, err
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