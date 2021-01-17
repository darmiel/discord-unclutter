package mayfly

import (
	"github.com/bwmarrin/discordgo"
	"log"
	"time"
)

type MayflyMessage struct {
	Message     *discordgo.Message
	DeleteAfter time.Duration
}

func NewMessage(msg *discordgo.Message, delay time.Duration) *MayflyMessage {
	if delay == 0 {
		delay = 10 * time.Second
	}
	return &MayflyMessage{
		Message:     msg,
		DeleteAfter: delay,
	}
}

func Queue(msg *discordgo.Message, delay time.Duration) {
	if msg == nil {
		log.Println("Skipped nil message on queue with duration:", delay)
	}
	message := NewMessage(msg, delay)
	mayflies.PushBack(message)
}

func QueueDefault(msg *discordgo.Message) {
	Queue(msg, 0)
}
