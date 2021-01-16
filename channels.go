package main

import (
	"errors"
	"github.com/bwmarrin/discordgo"
	"log"
	"strings"
)

// ChannelCategory is a tuple of a channel and the category
type ChannelCategory struct {
	Channel  *discordgo.Channel
	Category *discordgo.Channel
}

type UnclutteredChannel struct {
	*discordgo.Channel
	VoiceChannelID string
}

var channelCache = make(map[string]*ChannelCategory)

func findUnclutteredChannels(s *discordgo.Session, guildID string) (channels []*UnclutteredChannel, category *discordgo.Channel, err error) {
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
		channels, category, _ := findUnclutteredChannels(us.Session, us.GuildID)

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
		log.Println("🎉 Creating category", CategoryName)
		catChannel, err = us.Session.GuildChannelCreateComplex(us.GuildID, discordgo.GuildChannelCreateData{
			Name: CategoryName,
			Type: discordgo.ChannelTypeGuildCategory,
		})
		if err != nil {
			log.Println("ERROR: (creating category)", err)
			return nil, err
		}
	}

	if chChannel == nil {
		log.Println("🎉 Creating channel for voice", us.ChannelID)

		// create channel
		overwrites := []*discordgo.PermissionOverwrite{
			{
				ID:    us.GuildID, // Guild ID = everyone
				Type:  "0",
				Allow: 0,
				Deny:  discordgo.PermissionViewChannel, /* | discordgo.PermissionReadMessageHistory */
			},
			{
				ID:    us.Session.State.User.ID,
				Type:  "1",
				Allow: discordgo.PermissionAllChannel, // all access to channel
				Deny:  0,
			},
		}

		// get channel
		channel, err := us.Session.Channel(us.ChannelID)
		if err != nil || channel == nil {
			log.Println("ERROR: (retrieving channel info)", err)
			return nil, err
		}

		name := friendlyChannelName(channel.Name)
		log.Println("    🎉 Friendly:", name)

		chChannel, err = us.Session.GuildChannelCreateComplex(us.GuildID, discordgo.GuildChannelCreateData{
			Name:                 name,
			Type:                 discordgo.ChannelTypeGuildText,
			ParentID:             catChannel.ID,
			Topic:                TopicPrefix + us.ChannelID,
			PermissionOverwrites: overwrites,
		})
		if err != nil {
			log.Println("ERROR: (creating text channel)", err)
			return nil, err
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

func extractChannelID(channel *discordgo.Channel) (voiceID string) {
	if strings.HasPrefix(channel.Topic, TopicPrefix) {
		return strings.TrimSpace(channel.Topic[len(TopicPrefix):])
	}
	return ""
}
