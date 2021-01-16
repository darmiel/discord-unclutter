package unclutterer

import (
	"container/list"
	"github.com/bwmarrin/discordgo"
	"github.com/darmiel/discord-unclutterer/internal/unclutterer/database"
	"log"
	"strings"
	"time"
)

var messageCache = make(map[string]*discordgo.Message)
var messageNotif = list.New()

func DeleteNotifications(s *discordgo.Session, c chan bool) {
	ticker := time.NewTicker(1 * time.Second)

	for {
		select {
		case <-c:
			log.Println("Stop deleting notifications")
			return
		case <-ticker.C:
			for e := messageNotif.Front(); e != nil; e = e.Next() {
				msg, ok := e.Value.(*discordgo.Message)
				if !ok {
					log.Println("Error with element:", *e)
					continue
				}

				parse, err := msg.Timestamp.Parse()
				if err != nil {
					log.Println("Error parsing time on element:", *e)
					continue
				}

				if parse.Before(time.Now().Add(-10 * time.Second)) {
					// delete message
					log.Println("Delete message")

					if err := s.ChannelMessageDelete(msg.ChannelID, msg.ID); err != nil {
						log.Println("  └ ❌", err)
					}

					messageNotif.Remove(e)
				}
			}
		}
	}
}

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

func HandleMessageReactionAdd(s *discordgo.Session, ev *discordgo.MessageReactionAdd) {

	// Ignore self
	if ev.UserID == s.State.User.ID {
		return
	}

	if ev.Emoji.Name != Reaction {
		return
	}

	message := getMessage(s, ev.ChannelID, ev.MessageID)
	if message == nil {
		return
	}

	if !strings.HasPrefix(message.Content, ReactionCommand) {
		return
	}

	log.Println("has prefix!")

	msg, _ := s.ChannelMessageSend(ev.ChannelID, "[ <@"+ev.UserID+"> ] 👉 **Opt-Out** Ghost-Pings ... (loading)")
	if msg != nil {
		messageNotif.PushBack(msg)
	}

	// opt out
	if err := database.SetBlocksGhostping(ev.UserID, true); err != nil {
		log.Println("Error opt-out:", err)

		if msg != nil {
			_, _ = s.ChannelMessageEdit(ev.ChannelID, msg.ID, "[ <@"+ev.UserID+"> ] 👉 😡 **Nope:** "+err.Error())
		}
	} else {
		log.Println("Opt-out successful.")
		if msg != nil {
			_, _ = s.ChannelMessageEdit(
				ev.ChannelID,
				msg.ID,
				"[ <@"+ev.UserID+"> | https://tenor.com/wroQ.gif ] 👉 😊 Okay! Du erhältst keine weiteren Ghost-Pings",
			)
		}
	}
}

func HandleMessageReactionRemove(s *discordgo.Session, ev *discordgo.MessageReactionRemove) {

	// Ignore self
	if ev.UserID == s.State.User.ID {
		return
	}

	if ev.Emoji.Name != Reaction {
		return
	}

	message := getMessage(s, ev.ChannelID, ev.MessageID)
	if message == nil {
		return
	}

	if !strings.HasPrefix(message.Content, ReactionCommand) {
		return
	}

	msg, _ := s.ChannelMessageSend(ev.ChannelID, "[ <@"+ev.UserID+"> ] 👈 **Opt-In** Ghost-Pings ... (loading)")
	if msg != nil {
		messageNotif.PushBack(msg)
	}

	// opt out
	if err := database.SetBlocksGhostping(ev.UserID, false); err != nil {
		log.Println("Error opt-in:", err)

		if msg != nil {
			_, _ = s.ChannelMessageEdit(ev.ChannelID, msg.ID, "[ <@"+ev.UserID+"> ] 👈 😡 **Nope:** "+err.Error())
		}
	} else {
		log.Println("Opt-in successful.")
		if msg != nil {
			_, _ = s.ChannelMessageEdit(
				ev.ChannelID,
				msg.ID,
				"[ <@"+ev.UserID+"> | https://tenor.com/v4hv.gif ] 👈 😊 Okay! Du erhältst wieder Ghost-Pings",
			)
		}
	}
}
