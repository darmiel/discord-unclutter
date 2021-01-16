package unclutterer

import (
	"github.com/bwmarrin/discordgo"
	"log"
)

const (
	Reaction        = "🖕"
	ReactionCommand = "cmd::opt+in:out"
)

func (us *UserVoiceStateSession) CreateCategory() (channel *discordgo.Channel, err error) {
	log.Println("🎉🐈 Creating category", CategoryName)
	channel, err = us.Session.GuildChannelCreateComplex(us.GuildID, discordgo.GuildChannelCreateData{
		Name: CategoryName,
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
	}

	// get channel
	voiceChannel, err := us.Session.Channel(us.ChannelID)
	if err != nil || voiceChannel == nil {
		log.Println("ERROR: (retrieving channel info)", err)
		return nil, err
	}

	name := friendlyChannelName(voiceChannel.Name)
	log.Println("    └ 🎉 Friendly:", name)

	channel, err = us.Session.GuildChannelCreateComplex(us.GuildID, discordgo.GuildChannelCreateData{
		Name:                 name,
		Type:                 discordgo.ChannelTypeGuildText,
		ParentID:             parentID,
		Topic:                TopicPrefix + us.ChannelID,
		PermissionOverwrites: permissions,
	})

	log.Println("Channel:", channel, "err:", err)
	if channel != nil && err == nil {
		log.Print("  └ Creating welcome message")
		if _, err := us.SendWelcomeMessage(channel, voiceChannel); err != nil {
			log.Println("ERROR sending welcome message:", err)
		}
	}

	return
}

func (us *UserVoiceStateSession) SendWelcomeMessage(channel *discordgo.Channel, voiceChannel *discordgo.Channel) (message *discordgo.Message, err error) {
	var text = ReactionCommand + `
Hallo [ @everyone | https://i.imgur.com/aHX3n0z.png ]!

Dieser Channel wurde für den Voice-Channel ` + "`" + voiceChannel.Name + "`" + ` erstellt.
Er wird nur dann sichtbar, wenn du in diesen Voice-Channel gehst. (Privater Textkanal für Sprachkanäle).

Jedes Mal, wenn du Zugriff zu einem solchen Text-Channel bekommst, erhälst du einen Ghost-Ping.
Möchtest du diese Ghost-Pings nicht mehr erhalten, klicke auf '` + Reaction + `'`

	// send message
	log.Println("    └ Sending message")
	message, err = us.Session.ChannelMessageSend(channel.ID, text)
	if err != nil {
		return
	}

	// add middle finger reaction
	log.Println("    └ Adding reaction")
	err = us.Session.MessageReactionAdd(channel.ID, message.ID, Reaction)
	if err != nil {
		return
	}

	// pin message
	log.Println("    └ Pin")
	err = us.Session.ChannelMessagePin(channel.ID, message.ID)

	return
}
