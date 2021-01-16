package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/darmiel/discord-unclutterer/internal/unclutterer"
	"github.com/darmiel/discord-unclutterer/internal/unclutterer/cooldown"
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
		sess := &unclutterer.UserSess{
			UserID:    u.UserID,
			ChannelID: u.ChannelID,
			GuildID:   u.GuildID,
			Session:   s,
			Previous:  u.BeforeUpdate,
		}

		/*
			JOIN:
			2021/01/15 21:23:05 Current-State: &{0xc0003960a0 <nil>}
			2021/01/15 21:23:05   channel: 799732716340772913 user: 150347348088848384 guild: 779370167699898399
			2021/01/15 21:23:05 Previous-State: <nil>
			2021/01/15 21:23:05   (nothing)
		*/

		/*
			SWITCH:
			2021/01/15 21:23:23 Current-State: &{0xc000396140 0xc000396190}
			2021/01/15 21:23:23   channel: 799732735282380810 user: 150347348088848384 guild: 779370167699898399
			2021/01/15 21:23:23 Previous-State: &{150347348088848384 f666b9e60388a90ec7491e013c886e82 799732716340772913 779370167699898399 false true false false false}
			2021/01/15 21:23:23   channel: 799732716340772913 user: 150347348088848384 guild: 779370167699898399
		*/

		/*
			LEAVE:
			2021/01/15 21:23:45 Current-State: &{0xc000396370 0xc0003963c0}
			2021/01/15 21:23:45   channel:  user: 150347348088848384 guild: 779370167699898399
			2021/01/15 21:23:45 Previous-State: &{150347348088848384 f666b9e60388a90ec7491e013c886e82 799732735282380810 779370167699898399 false true false false false}
			2021/01/15 21:23:45   channel: 799732735282380810 user: 150347348088848384 guild: 779370167699898399
		*/

		if u == nil {
			return
		}

		// check leave
		if u.ChannelID == "" || u.BeforeUpdate != nil {
			sess.UserLeave()
			if u.ChannelID == "" || (u.BeforeUpdate != nil && u.BeforeUpdate.ChannelID == "") {
				return
			}
		}

		// check cooldown
		if cd, vl := cooldown.IsOnCooldown(u.UserID); cd {
			log.Println("  ⏰  User", u.UserID, "on cooldown! ( VL:", vl, ")")
			if vl >= 3 {
				log.Println("WARN: User", u.UserID, "has a high amount of violations!")
			}
			return
		}

		sess.UserJoin()
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
		channels, _, _ := unclutterer.FindCreatedTextChannels(discord, guild.ID)
		if channels == nil {
			log.Println("Error checking guild:", guild.ID)
			continue
		}

		for _, channel := range channels {
			log.Println("  └ Channel:", channel.ID, "(", channel.Name, ")")

			var cleared = false
			for _, perm := range channel.PermissionOverwrites {
				// ignore @everyone and self
				if perm.ID == guild.ID || perm.ID == discord.State.User.ID {
					continue
				}
				// remove
				if err := unclutterer.RevokeAccess(discord, perm.ID, channel.ID, false); err != nil {
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

	discord.AddHandler(func(s *discordgo.Session, e *discordgo.MessageCreate) {
		chanId := e.ChannelID
		channel, err := s.State.Channel(chanId)
		if err != nil {
			log.Println("Error receiving channel:", err)
			return
		}
		fmt.Println("CHAT |", channel.Name, "(", e.Author.Username, "):", e.Content)
	})

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	log.Println("Close:", discord.Close())
}
