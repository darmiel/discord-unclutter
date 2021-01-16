package mayfly

import (
	"container/list"
	"github.com/bwmarrin/discordgo"
	"log"
	"strings"
	"time"
)

var mayflies = list.New()

func DeleteNotifications(session *discordgo.Session, cancel chan bool) {
	ticker := time.NewTicker(1 * time.Second)

	for {
		select {

		case <-cancel:
			log.Println("Stop deleting notifications")
			return

		case <-ticker.C:
			for e := mayflies.Front(); e != nil; e = e.Next() {
				mf, ok := e.Value.(*MayflyMessage)
				if !ok {
					log.Println("Error with element:", *e)
					continue
				}

				// get discord message
				msg := mf.Message

				// parse timestamp
				parse, err := msg.Timestamp.Parse()
				if err != nil {
					log.Println("Error parsing time on element:", *e)
					continue
				}

				// get delete after time
				if mf.DeleteAfter == 0 {
					mf.DeleteAfter = 10 * time.Second
				}

				// parse
				if parse.Before(time.Now().Add(-1 * mf.DeleteAfter)) {
					// build content preview
					content := strings.ReplaceAll(msg.Content, "\n", "\\n")
					if len(content) > 64 {
						content = content[:60] + " ..."
					}
					//

					// delete message
					log.Println("🪰 Deleting mayfly message:", msg.ID, "(", content, ")")
					if err := session.ChannelMessageDelete(msg.ChannelID, msg.ID); err != nil {
						log.Println("  └ ❌", err)
					} else {
						log.Println("  └ ✅ Message deleted.")
					}

					mayflies.Remove(e)
				}
			}
		}
	}
}
