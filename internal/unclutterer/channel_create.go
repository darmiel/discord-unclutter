package unclutterer

import (
	"github.com/bwmarrin/discordgo"
	"log"
)

func (us *UserSess) CreateCategory() (channel *discordgo.Channel, err error) {
	log.Println("🎉🐈 Creating category", CategoryName)
	channel, err = us.Session.GuildChannelCreateComplex(us.GuildID, discordgo.GuildChannelCreateData{
		Name: CategoryName,
		Type: discordgo.ChannelTypeGuildCategory,
	})
	return
}

func (us *UserSess) CreateChannel(parentID string) (channel *discordgo.Channel, err error) {
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

	return
}
