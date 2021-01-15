package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
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
	log.Println("Token:", discordToken)

	discord, err := discordgo.New("Bot " + discordToken)
	if err != nil {
		log.Fatalln("Error creating discord:", err)
		return
	}

	discord.AddHandler(func(s *discordgo.Session, u *discordgo.VoiceStateUpdate) {
		sess := &UserSess{
			UserID:    u.UserID,
			ChannelID: u.ChannelID,
			GuildID:   u.GuildID,
			Session:   s,
			Previous:  u.BeforeUpdate,
		}
		leftChannel := u.BeforeUpdate != nil

		if leftChannel {
			sess.userLeaveChannel()
		} else {
			sess.userJoinChannel()
		}
	})

	err = discord.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// check all channels
	// iterate through every guild
	log.Println("┌ Starting guild channel clear")
	for _, guild := range discord.State.Guilds {
		log.Println("└", guild.ID, "(", guild.Name, ")")

		// find channels
		channels, _, _ := findUnclutteredChannels(discord, guild.ID)
		if channels == nil {
			log.Println("Error checking guild:", guild.ID)
			continue
		}

		for _, channel := range channels {
			log.Println("  └ Channel:", channel.ID, "(", channel.Name, ")")

			var cleared = false
			for _, perm := range channel.PermissionOverwrites {
				// ignore @everyone
				if perm.ID == guild.ID {
					continue
				}

				// remove
				if err := revokeAccess(discord, perm.ID, channel.ID, false); err != nil {
					log.Println("    └ ❌", err)
				} else {
					log.Println("    └ ✅", perm.ID)
					cleared = true
				}
			}

			if cleared {
				_, _ = discord.ChannelMessageSend(channel.ID, "🗑 Cleared all user's permissions")
			}
		}
	}

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	log.Println("Close:", discord.Close())
}
