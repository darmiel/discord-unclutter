package unclutterer

import (
	"github.com/bwmarrin/discordgo"
	duconfig "github.com/darmiel/discord-unclutterer/internal/unclutterer/config"
	"github.com/darmiel/discord-unclutterer/internal/unclutterer/cooldown"
	"log"
)

var selfAttributes = make(map[string]string)

func HandleVoiceStateUpdate(s *discordgo.Session, u *discordgo.VoiceStateUpdate, config *duconfig.Config) {
	if u == nil {
		return
	}

	// small & ugly hack to ignore mutes / un-mutes, deaves / un-deaves
	if previousChannel, ok := selfAttributes[u.UserID]; ok {
		if previousChannel == u.ChannelID {
			if config.VerbosityLevel >= 2 {
				log.Println("🟣 Ignored because same channel")
			}
			return
		}
	}
	selfAttributes[u.UserID] = u.ChannelID

	sess := &UserVoiceStateSession{
		UserID:    u.UserID,
		ChannelID: u.ChannelID,
		GuildID:   u.GuildID,
		Session:   s,
		Previous:  u.BeforeUpdate,
		config:    config,
	}

	// check leave
	if u.ChannelID == "" || u.BeforeUpdate != nil {
		// execute user leave
		sess.UserLeave(config)
		if u.ChannelID == "" || (u.BeforeUpdate != nil && u.BeforeUpdate.ChannelID == "") {
			return
		}
	}

	// check cool down
	if cd, vl := cooldown.IsOnCooldown(u.UserID, config.JoinLeaveCooldown.Duration, config); cd {
		if config.LogCooldown {
			log.Println("⏰  User", u.UserID, "on cool down! ( VL:", vl, ")")
		}
		return
	}

	// execute user join
	sess.UserJoin(config)
}
