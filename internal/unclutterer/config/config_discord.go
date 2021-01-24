package config

type DiscordConfig struct {
	DiscordToken        string
	ChannelTopicPrefix  string
	ChannelCategoryName string
	ChannelTextPrefix   string

	CleanChannelsOnStartup bool
	MayflyCheckInterval    duration

	DiscordGameStatus string
}
