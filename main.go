package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/darmiel/discord-unclutterer/internal/unclutterer"
	"github.com/darmiel/discord-unclutterer/internal/unclutterer/cleanup"
	duconfig "github.com/darmiel/discord-unclutterer/internal/unclutterer/config"
	"github.com/darmiel/discord-unclutterer/internal/unclutterer/mayfly"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	// load config
	config, err := duconfig.LoadConfig()
	if err != nil {
		log.Fatalln("Error decoding config:", err)
		return
	}

	// check token
	discordToken := config.DiscordToken
	if discordToken == "" {
		log.Fatalln("❌ No Discord token specified.")
		return
	}
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
	if config.CleanChannelsOnStartup {
		log.Println("🗑 Starting guild channel clear")
		log.Println("   └ To disable CleanChannelsOnStartup in Knotig set to False.")

		cleanup.StartCleanup(discord, config)
	}

	// register handlers
	discord.AddHandler(func(s *discordgo.Session, e *discordgo.MessageCreate) {
		unclutterer.HandleMessageCreate(s, e, config)
	})
	discord.AddHandler(func(s *discordgo.Session, e *discordgo.VoiceStateUpdate) {
		unclutterer.HandleVoiceStateUpdate(s, e, config)
	})
	discord.AddHandler(func(s *discordgo.Session, e *discordgo.MessageReactionAdd) {
		unclutterer.HandleMessageReactionAdd(s, e, config)
	})
	discord.AddHandler(func(s *discordgo.Session, e *discordgo.MessageReactionRemove) {
		unclutterer.HandleMessageReactionRemove(s, e, config)
	})

	// notification deleter
	done := make(chan bool)
	go func() {
		log.Println("🪰 Starting mayfly task with a interval of", config.MayflyCheckInterval.String())
		mayfly.DeleteNotifications(discord, config, done)
	}()
	//

	// update activity
	if config.DiscordGameStatus != "" {
		log.Println("📝 Updating status to:", config.DiscordGameStatus)
		idle := int(time.Now().UnixNano() / int64(1000000)) // convert unix nano to unix ms
		if err := discord.UpdateStatus(idle, config.DiscordGameStatus); err != nil {
			log.Println("  └ ❌ Error updating status:", err)
		} else {
			log.Println("  └ ✅  Status updated!")
		}
	}

	fmt.Println("")
	fmt.Println("+-------------------------------------------+")
	fmt.Println("| Bot is now running. Press CTRL-C to exit. |")
	fmt.Println("+-------------------------------------------+")
	fmt.Println("")

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
