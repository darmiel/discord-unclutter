package config

type DiscordConfig struct {
	DiscordToken        string
	ChannelTopicPrefix  string
	ChannelCategoryName string
	ChannelTextPrefix   string

	OptReaction        string
	OptReactionCommand string

	CleanChannelsOnStartup bool
	MayflyCheckInterval    duration

	DiscordGameStatus string
}
