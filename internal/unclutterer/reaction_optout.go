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
	opt(s, ev.MessageReaction, true, config)
}

func HandleMessageReactionRemove(s *discordgo.Session, ev *discordgo.MessageReactionRemove, config *duconfig.Config) {
	opt(s, ev.MessageReaction, false, config)
}

func opt(s *discordgo.Session, reaction *discordgo.MessageReaction, block bool, config *duconfig.Config) {

	// Ignore self
	if reaction.UserID == s.State.User.ID {
		return
	}
	if reaction.Emoji.Name != config.OptReaction {
		return
	}
	// check cool down
	if CheckCooldown(reaction.UserID, config) {
		return
	}
	message := getMessage(s, reaction.ChannelID, reaction.MessageID)
	if message == nil {
		return
	}
	if !strings.HasPrefix(message.Content, config.OptReactionCommand) {
		return
	}

	var msg *discordgo.Message

	if block {
		msg, _ = s.ChannelMessageSend(
			reaction.ChannelID,
			config.OptOutLoading.Repl("UserID", reaction.UserID),
		)
	} else {
		msg, _ = s.ChannelMessageSend(
			reaction.ChannelID,
			config.OptInLoading.Repl("UserID", reaction.UserID),
		)
	}
	mayfly.QueueDefault(msg)

	if !config.AllowGhostPingBlocking {
		_, _ = s.ChannelMessageEdit(
			reaction.ChannelID,
			msg.ID,
			config.OptDisabled.Repl("UserID", reaction.UserID),
		)
		return
	}

	// opt in/out
	if err := database.SetBlocksGhostping(reaction.UserID, block, config); err != nil {
		if config.LogOptOutError {
			if block {
				log.Println("❌👉 Error opting-out for", reaction.UserID, ":", err)
			} else {
				log.Println("❌👉 Error opting-in for", reaction.UserID, ":", err)
			}
		}

		if msg != nil {
			if block {
				_, _ = s.ChannelMessageEdit(
					reaction.ChannelID,
					msg.ID,
					config.OptOutErrorData.Repl("UserID", reaction.UserID, "Error", err.Error()),
				)
			} else {
				_, _ = s.ChannelMessageEdit(
					reaction.ChannelID,
					msg.ID,
					config.OptInErrorDatabase.Repl("UserID", reaction.UserID, "Error", err.Error()),
				)
			}
		}
	} else {
		if config.LogOptOutSuccess {
			if block {
				log.Println("✅ 👉 Opt-out for", reaction.UserID, "successful.")
			} else {
				log.Println("✅ 👉 Opt-in for", reaction.UserID, "successful.")
			}
		}

		if msg != nil {
			if block {
				_, _ = s.ChannelMessageEdit(
					reaction.ChannelID,
					msg.ID,
					config.OptOutSuccess.Repl("UserID", reaction.UserID),
				)
			} else {
				_, _ = s.ChannelMessageEdit(
					reaction.ChannelID,
					msg.ID,
					config.OptInSuccess.Repl("UserID", reaction.UserID),
				)
			}
		}
	}
}
