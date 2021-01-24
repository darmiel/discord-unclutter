package unclutterer

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	duconfig "github.com/darmiel/discord-unclutterer/internal/unclutterer/config"
	"log"
)

func HandleMessageCreate(s *discordgo.Session, e *discordgo.MessageCreate, config *duconfig.Config) {
	chanId := e.ChannelID
	channel, err := s.State.Channel(chanId)
	if err != nil {
		if config.VerbosityLevel >= 1 {
			log.Println("❌ Error receiving channel on message create:", err)
		}
		return
	}
	if config.LogChat {
		fmt.Println("CHAT |", channel.Name, "(", e.Author.Username, "):", e.Content)
	}
}
