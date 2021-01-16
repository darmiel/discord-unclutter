package main

import (
	"github.com/bwmarrin/discordgo"
	"log"
	"regexp"
)

var (
	channelNameRegex *regexp.Regexp
)

const (
	TopicPrefix  = "dcuncltr: "
	CategoryName = "VOICE TEXT CHANNELS"
)

func init() {
	channelNameRegex = regexp.MustCompile("[^\\w-_]]")
}

func friendlyChannelName(channel string) (res string) {
	res = channelNameRegex.ReplaceAllString(channel, "")
	return
}

type UserSess struct {
	UserID    string
	ChannelID string
	GuildID   string
	Session   *discordgo.Session
	Previous  *discordgo.VoiceState
}

func (us *UserSess) userJoinChannel() {
	log.Println("User", us.UserID, "joined", "channel", us.ChannelID, "from guild", us.GuildID)

	// find channel
	if channelPair, err := us.findOrCreateText(us.ChannelID); err != nil {
		log.Println("ERROR on finding text channel:", err)
		return
	} else {
		textChannel := channelPair.Channel

		if err := grantAccess(us.Session, us.UserID, textChannel.ID); err != nil {
			log.Println("ERROR: Granting Access (", textChannel.Name, ") for", us.UserID, ":", err)
			return
		}

		// ping user
		if _, err := us.Session.ChannelMessageSend(
			textChannel.ID,
			"✅ <@"+us.UserID+"> wurde hinzugefügt 👋",
		); err != nil {
			log.Println("ERROR sending message:", err)
		}
	}
}

func (us *UserSess) userLeaveChannel() {
	if us.Previous == nil {
		log.Println("User", us.UserID, "left unknown channel.")
		return
	}

	log.Println("User", us.UserID, "left", "channel", us.Previous.ChannelID, "from guild", us.GuildID)

	// find channel
	if channelPair, err := us.findOrCreateText(us.Previous.ChannelID); err != nil {
		log.Println("ERROR on finding text channel:", err)
		return
	} else {
		textChannel := channelPair.Channel

		if err := revokeAccess(us.Session, us.UserID, textChannel.ID, false); err != nil {
			log.Println("ERROR: Revoking Access (", textChannel.Name, ") for", us.UserID, ":", err)
			return
		}

		// ping user
		if _, err := us.Session.ChannelMessageSend(
			textChannel.ID,
			"🚪 <@"+us.UserID+"> wurde entfernt 👋",
		); err != nil {
			log.Println("ERROR sending message:", err)
		}
	}
}
