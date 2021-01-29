package unclutterer

import (
	duconfig "github.com/darmiel/discord-unclutterer/internal/unclutterer/config"
	"github.com/darmiel/discord-unclutterer/internal/unclutterer/database"
	"log"
	"time"
)

func (us *UserVoiceStateSession) MentionUser() string {
	return "<@" + us.UserID + ">"
}

func (us *UserVoiceStateSession) UserJoin(config *duconfig.Config) {
	if config.LogUserJoin {
		log.Println("👋", us.UserID, "joined", us.ChannelID)
	}

	// find channel
	channelPair, err := us.findOrCreateText(us.ChannelID)
	if err != nil {
		// channel not found
		log.Println("ERROR on finding text channel:", err)
		return
	}
	// channel found

	textChannel := channelPair.Channel

	if err := GrantAccess(us.Session, us.UserID, textChannel.ID); err != nil {
		log.Println("ERROR: Granting Access (", textChannel.Name, ") for", us.UserID, ":", err)
		return
	}

	t := time.Now()

	block := config.GhostPingBlockDefault
	if config.AllowGhostPingBlocking {
		// check if user wants to receive ghost pings (opt-out)
		block, err = database.BlocksGhostping(us.UserID, config, block)
		if config.VerbosityLevel >= 3 {
			log.Println("👻 Ghost-Ping Get Result:", block, err)
		}
		if err != nil {
			if config.VerbosityLevel >= 3 {
				log.Println("💾 Database (Ghost-Ping) Error:", err)
			}
			return
		}
		if config.VerbosityLevel >= 3 {
			log.Println("   └ Get Blocks took", time.Now().Unix()-t.Unix(), "s")
		}
	} else {
		if config.VerbosityLevel >= 3 {
			log.Println("👻 Ghost-Ping:", block, "(default, config)")
		}
	}

	// make ghost ping
	if !block {
		send, err := us.Session.ChannelMessageSend(textChannel.ID, us.MentionUser())
		if err != nil {
			if config.VerbosityLevel >= 3 {
				log.Println("   └ 👻 Ghost-Ping Create Error:", err)
			}
			return
		}

		if err := us.Session.ChannelMessageDelete(textChannel.ID, send.ID); err != nil {
			if config.VerbosityLevel >= 3 {
				log.Println("   └ 👻 Ghost-Ping Delete Error:", err)
			}
		}
	}
}

func (us *UserVoiceStateSession) UserLeave(config *duconfig.Config) {
	if config.LogUserLeave {
		log.Println("🚪", us.UserID, "left", us.ChannelID)
	}

	if us.Previous == nil {
		return
	}

	// find channel
	if channelPair, err := us.findOrCreateText(us.Previous.ChannelID); err != nil {
		log.Println("ERROR on finding text channel:", err)
		return
	} else {
		textChannel := channelPair.Channel

		if err := RevokeAccess(us.Session, us.UserID, textChannel.ID, false); err != nil {
			log.Println("ERROR: Revoking Access (", textChannel.Name, ") for", us.UserID, ":", err)
			return
		}
	}
}
