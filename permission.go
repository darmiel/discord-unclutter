package main

import (
	"github.com/bwmarrin/discordgo"
	"log"
)

func grantAccess(s *discordgo.Session, userID string, channelID string) (err error) {
	if userID == s.State.User.ID {
		log.Println("Skipped grant access for", userID, "because it was me!")
		return nil
	}

	return s.ChannelPermissionSet(
		channelID,
		userID,
		"1",
		discordgo.PermissionViewChannel,
		0,
	)
}

func revokeAccess(s *discordgo.Session, userID string, channelID string, force bool) (err error) {
	if userID == s.State.User.ID {
		log.Println("Skipped revoke access for", userID, "because it was me!")
		return nil
	}
	if force {
		return s.ChannelPermissionSet(
			channelID,
			userID,
			"1",
			0,
			discordgo.PermissionViewChannel,
		)
	}

	return s.ChannelPermissionDelete(
		channelID,
		userID,
	)
}
