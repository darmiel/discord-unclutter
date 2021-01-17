package unclutterer

import (
	"github.com/bwmarrin/discordgo"
	duconfig "github.com/darmiel/discord-unclutterer/internal/unclutterer/config"
	"github.com/darmiel/discord-unclutterer/internal/unclutterer/cooldown"
	"github.com/darmiel/discord-unclutterer/internal/unclutterer/database"
	"github.com/darmiel/discord-unclutterer/internal/unclutterer/mayfly"
	"log"
	"strings"
)

var messageCache = make(map[string]*discordgo.Message)

func getMessage(s *discordgo.Session, channelID string, messageID string) (message *discordgo.Message) {
	log.Println("☎️ Get message", messageID)
	if msg, ok := messageCache[messageID]; ok {
		log.Println("  └ From Cache")
		return msg
	}

	log.Println("  └ From Discord")

	// get

	message, _ = s.ChannelMessage(channelID, messageID)
	if message != nil {
		messageCache[messageID] = message
	}

	return
}

func CheckCooldown(u string, config *duconfig.Config) bool {
	c, _ := cooldown.IsOnCooldown(u+":opt::in+out", config.OptCooldown.Duration, config)
	return c
}

func HandleMessageReactionAdd(s *discordgo.Session, ev *discordgo.MessageReactionAdd, config *duconfig.Config) {
	// Ignore self
	if ev.UserID == s.State.User.ID {
		return
	}

	if ev.Emoji.Name != config.OptReaction {
		return
	}

	// check cool down
	if CheckCooldown(ev.UserID, config) {
		return
	}

	message := getMessage(s, ev.ChannelID, ev.MessageID)
	if message == nil {
		return
	}

	if !strings.HasPrefix(message.Content, config.OptReactionCommand) {
		return
	}

	msg, _ := s.ChannelMessageSend(ev.ChannelID, config.OptOutLoading.Repl("UserID", ev.UserID))
	mayfly.QueueDefault(msg)

	if !config.AllowGhostPingBlocking {
		_, _ = s.ChannelMessageEdit(
			ev.ChannelID,
			msg.ID,
			config.OptDisabled.Repl("UserID", ev.UserID),
		)
		return
	}

	// opt out
	if err := database.SetBlocksGhostping(ev.UserID, true, config); err != nil {
		if config.LogOptOutError {
			log.Println("❌👉 Error opting-out for", ev.UserID, ":", err)
		}

		if msg != nil {
			_, _ = s.ChannelMessageEdit(
				ev.ChannelID,
				msg.ID,
				config.OptOutErrorData.Repl("UserID", ev.UserID, "Error", err.Error()),
			)
		}
	} else {
		if config.LogOptOutSuccess {
			log.Println("✅ 👉 Opt-out for", ev.UserID, "successful.")
		}

		if msg != nil {
			_, _ = s.ChannelMessageEdit(
				ev.ChannelID,
				msg.ID,
				config.OptOutSuccess.Repl("UserID", ev.UserID),
			)
		}
	}
}

func HandleMessageReactionRemove(s *discordgo.Session, ev *discordgo.MessageReactionRemove, config *duconfig.Config) {
	// Ignore self
	if ev.UserID == s.State.User.ID {
		return
	}

	if ev.Emoji.Name != config.OptReaction {
		return
	}

	// check cool down
	if CheckCooldown(ev.UserID, config) {
		return
	}

	message := getMessage(s, ev.ChannelID, ev.MessageID)
	if message == nil {
		return
	}

	if !strings.HasPrefix(message.Content, config.OptReactionCommand) {
		return
	}

	msg, _ := s.ChannelMessageSend(ev.ChannelID, config.OptInLoading.Repl("UserID", ev.UserID))
	mayfly.QueueDefault(msg)

	if !config.AllowGhostPingBlocking {
		_, _ = s.ChannelMessageEdit(
			ev.ChannelID,
			msg.ID,
			config.OptDisabled.Repl("UserID", ev.UserID),
		)
		return
	}

	// opt in
	if err := database.SetBlocksGhostping(ev.UserID, false, config); err != nil {
		if config.LogOptInError {
			log.Println("❌👈 Error opting-in for", ev.UserID, ":", err)
		}

		if msg != nil {
			_, _ = s.ChannelMessageEdit(
				ev.ChannelID,
				msg.ID,
				config.OptInErrorDatabase.Repl("UserID", ev.UserID, "Error", err.Error()),
			)
		}
	} else {
		if config.LogOptInSuccess {
			log.Println("✅ 👈 Opt-in for", ev.UserID, "successful.")
		}

		if msg != nil {
			_, _ = s.ChannelMessageEdit(
				ev.ChannelID,
				msg.ID,
				config.OptInSuccess.Repl("UserID", ev.UserID),
			)
		}
	}
}
