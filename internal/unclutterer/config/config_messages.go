package config

type MessageConfig struct {
	ChannelCreateMessage PlaceholderMessage

	OptInLoading       PlaceholderMessage
	OptInErrorDatabase PlaceholderMessage
	OptInSuccess       PlaceholderMessage

	OptOutLoading   PlaceholderMessage
	OptOutErrorData PlaceholderMessage
	OptOutSuccess   PlaceholderMessage

	OptDisabled PlaceholderMessage
}
