package unclutterer

import "github.com/bwmarrin/discordgo"

// ChannelCategory is a tuple of a channel and the category
type ChannelCategory struct {
	Channel  *discordgo.Channel
	Category *discordgo.Channel
}

type UnclutteredChannel struct {
	*discordgo.Channel
	VoiceChannelID string
}

type UserVoiceStateSession struct {
	UserID    string
	ChannelID string
	GuildID   string
	Session   *discordgo.Session
	Previous  *discordgo.VoiceState
}
