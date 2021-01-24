package cleanup

import (
	"github.com/bwmarrin/discordgo"
	"github.com/darmiel/discord-unclutterer/internal/unclutterer"
	duconfig "github.com/darmiel/discord-unclutterer/internal/unclutterer/config"
	"log"
)

// StartCleanup checks all text channels if they have any users' permissions left
// and removes them if necessary
func StartCleanup(s *discordgo.Session, config *duconfig.Config) {
	guilds := s.State.Guilds
	for gi, guild := range guilds {

		// build structure prefix
		var prefix string
		if (gi + 1) >= len(guilds) {
			prefix = "   └"
		} else {
			prefix = "   ├"
		}
		//

		log.Println(prefix, guild.ID, "(", guild.Name, ")")

		// find channels
		channels, _, _ := unclutterer.FindCreatedTextChannels(s, guild.ID, config)
		if channels == nil {
			log.Println("Error checking guild:", guild.ID)
			continue
		}

		for ci, channel := range channels {

			// build structure prefix
			if (ci + 1) >= len(channels) {
				prefix = "     └"
			} else {
				prefix = "     ├"
			}
			//

			log.Println(prefix, "Channel:", channel.ID, "(", channel.Name, ")")

			var cleared = false
			for _, perm := range channel.PermissionOverwrites {
				// ignore @everyone and self
				if perm.ID == guild.ID || perm.ID == s.State.User.ID || perm.ID == unclutterer.TemporaryFixID {
					continue
				}

				// remove
				if err := unclutterer.RevokeAccess(s, perm.ID, channel.ID, false); err != nil {
					log.Println("       ├ ❌", err)
				} else {
					log.Println("       ├ ✅ ", perm.ID)
					cleared = true
				}
			}

			if cleared {
				log.Println("       └ (Some) users were cleared")
			} else {
				log.Println("       └ No users were cleared")
			}
		}
	}
}
