package unclutterer

import (
	"errors"
	"github.com/bwmarrin/discordgo"
)

const (
	GrantPermission  = discordgo.PermissionViewChannel
	RevokePermission = discordgo.PermissionViewChannel

	// TODO: Temporary fix, remove later
	TemporaryFixID = "710491120903127080"
)

func GrantAccess(s *discordgo.Session, userID string, channelID string) (err error) {
	if userID == s.State.User.ID {
		return errors.New("tried to remove permission from me")
	}
	if userID == TemporaryFixID {
		return nil
	}

	return s.ChannelPermissionSet(
		channelID,
		userID,
		"1",
		GrantPermission,
		0,
	)
}

// force = negate permissions instead of deleting them
func RevokeAccess(s *discordgo.Session, userID string, channelID string, force bool) (err error) {
	if userID == s.State.User.ID {
		return errors.New("tried to remove permission from me")
	}
	if userID == TemporaryFixID {
		return nil
	}

	if force {
		return s.ChannelPermissionSet(
			channelID,
			userID,
			"1",
			0,
			RevokePermission,
		)
	}
	return s.ChannelPermissionDelete(
		channelID,
		userID,
	)
}
