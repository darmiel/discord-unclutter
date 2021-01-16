package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/darmiel/discord-unclutterer/internal/unclutterer"
	"github.com/darmiel/discord-unclutterer/internal/unclutterer/cleanup"
	"github.com/darmiel/discord-unclutterer/internal/unclutterer/mayfly"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var (
	discordToken string
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalln("Error loading .env file")
		return
	}

	discordToken = os.Getenv("DISCORD_TOKEN")
}

func main() {
	var tokenSecret string
	for i := 0; i < len(discordToken); i++ {
		tokenSecret += "*"
	}
	log.Println("🔑 Logging in with token", tokenSecret, "...")

	discord, err := discordgo.New("Bot " + discordToken)
	if err != nil {
		log.Fatalln("Error creating discord:", err)
		return
	}

	err = discord.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// start cleanup
	cleanup.StartCleanup(discord)

	// register handlers
	discord.AddHandler(unclutterer.HandleMessageCreate)
	discord.AddHandler(unclutterer.HandleVoiceStateUpdate)
	discord.AddHandler(unclutterer.HandleMessageReactionAdd)
	discord.AddHandler(unclutterer.HandleMessageReactionRemove)

	// notification deleter
	done := make(chan bool)
	go func() {
		mayfly.DeleteNotifications(discord, done)
	}()
	//

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// cancel delete
	done <- true

	if err := discord.Close(); err != nil {
		log.Println("❌ Discord (Session) could not be closed:", err)
	} else {
		log.Println("✅  Discord (Session) closed.")
	}
}
