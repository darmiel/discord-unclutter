package unclutterer

import (
	"github.com/bwmarrin/discordgo"
	"log"
)

func (us *UserVoiceStateSession) CreateCategory() (channel *discordgo.Channel, err error) {
	name := us.config.ChannelCategoryName

	log.Println("🎉🐈 Creating category", name)

	channel, err = us.Session.GuildChannelCreateComplex(us.GuildID, discordgo.GuildChannelCreateData{
		Name: name,
		Type: discordgo.ChannelTypeGuildCategory,
	})
	return
}

func (us *UserVoiceStateSession) CreateChannel(parentID string) (channel *discordgo.Channel, err error) {
	log.Println("🎉📺 Creating channel for voice", us.ChannelID)

	// permissions
	permissions := []*discordgo.PermissionOverwrite{
		{
			ID:    us.GuildID, // Guild ID = everyone
			Type:  "0",
			Allow: 0,
			Deny:  discordgo.PermissionViewChannel, /* | discordgo.PermissionReadMessageHistory */
		},
		{
			ID:    us.Session.State.User.ID, // bot
			Type:  "1",
			Allow: discordgo.PermissionAllChannel, // all access to channel
			Deny:  0,
		},

		// temporary fix -> remove this later
		{
			ID:    TemporaryFixID, // me
			Type:  "1",
			Allow: discordgo.PermissionAllChannel, // all access to channel
			Deny:  0,
		},
	}

	// get channel
	voiceChannel, err := us.Session.Channel(us.ChannelID)
	if err != nil || voiceChannel == nil {
		if us.config.VerbosityLevel >= 1 {
			log.Println("ERROR: (retrieving channel info)", err)
		}
		return nil, err
	}

	name := us.config.ChannelTextPrefix + friendlyChannelName(voiceChannel.Name)
	topic := us.config.ChannelTopicPrefix + us.ChannelID

	log.Println("    ├ 🎉 Name:", name)
	log.Println("    └ 🎉 Topic:", topic)

	channel, err = us.Session.GuildChannelCreateComplex(us.GuildID, discordgo.GuildChannelCreateData{
		Name:                 name,
		Type:                 discordgo.ChannelTypeGuildText,
		ParentID:             parentID,
		Topic:                topic,
		PermissionOverwrites: permissions,
	})

	if channel != nil && err == nil {
		if _, err := us.SendWelcomeMessage(channel, voiceChannel); err != nil {
			if us.config.VerbosityLevel >= 1 {
				log.Println("❌ Error sending welcome message:", err)
			}
		}
	}

	return
}

func (us *UserVoiceStateSession) SendWelcomeMessage(channel *discordgo.Channel, voiceChannel *discordgo.Channel) (message *discordgo.Message, err error) {
	text := us.config.ChannelCreateMessage.Repl(
		"UserID", us.UserID,
		"ChannelID", us.ChannelID,
		"GuildID", us.GuildID,
		"ReactionCommand", us.config.OptReactionCommand,
		"OptReaction", us.config.OptReaction,
		"VoiceChannelName", voiceChannel.Name,
		"VoiceChannelID", voiceChannel.ID,
	)

	// send message
	log.Println("    └ 💌 Sending message")
	message, err = us.Session.ChannelMessageSend(channel.ID, text)
	if err != nil {
		return
	}

	// add middle finger reaction
	log.Println("    └ 👍 Adding reaction")
	err = us.Session.MessageReactionAdd(channel.ID, message.ID, us.config.OptReaction)
	if err != nil {
		return
	}

	// pin message
	log.Println("    └ 📌 Pin")
	err = us.Session.ChannelMessagePin(channel.ID, message.ID)

	return
}
