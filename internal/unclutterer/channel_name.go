package unclutterer

import (
	"github.com/bwmarrin/discordgo"
	"regexp"
	"strings"
)

var channelNameRegex *regexp.Regexp

func init() {
	channelNameRegex = regexp.MustCompile("[^\\w-_]]")
}

const (
	// TODO: Use from config
	TopicPrefix = "dcuncltr: "
	// TODO: Use from config
	CategoryName = "VOICE TEXT CHANNELS"
)

func extractChannelID(channel *discordgo.Channel) (voiceID string) {
	if strings.HasPrefix(channel.Topic, TopicPrefix) {
		return strings.TrimSpace(channel.Topic[len(TopicPrefix):])
	}
	return ""
}

func friendlyChannelName(channel string) (res string) {
	res = channelNameRegex.ReplaceAllString(channel, "")
	return
}
