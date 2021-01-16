package unclutterer

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
)

func HandleMessageCreate(s *discordgo.Session, e *discordgo.MessageCreate) {
	chanId := e.ChannelID
	channel, err := s.State.Channel(chanId)
	if err != nil {
		log.Println("Error receiving channel:", err)
		return
	}
	// TODO: Use from config
	fmt.Println("CHAT |", channel.Name, "(", e.Author.Username, "):", e.Content)
}
