package unclutterer

import (
	"github.com/bwmarrin/discordgo"
	duconfig "github.com/darmiel/discord-unclutterer/internal/unclutterer/config"
	"regexp"
	"strings"
)

var channelNameRegex *regexp.Regexp

func init() {
	channelNameRegex = regexp.MustCompile("[^\\w-_]]")
}

func extractChannelID(channel *discordgo.Channel, config *duconfig.Config) (voiceID string) {
	if strings.HasPrefix(channel.Topic, config.ChannelTopicPrefix) {
		return strings.TrimSpace(channel.Topic[len(config.ChannelTopicPrefix):])
	}
	return ""
}

func friendlyChannelName(channel string) (res string) {
	res = channelNameRegex.ReplaceAllString(channel, "")
	return
}
