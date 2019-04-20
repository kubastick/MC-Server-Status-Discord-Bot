package main

import (
	"github.com/TurtleGamingFTW/dblgo-archive"
	"github.com/bwmarrin/discordgo"
	"log"
	"time"
)

func postServerCountToDiscordBotApi(s *discordgo.Session) {
	if config.DiscordStatsSecret != "" && config.DiscordAppID != "" {
		dblApi := dblgo.NewDBL(config.DiscordSecret, config.DiscordAppID)
		postStatsLoop(s, *dblApi)
	} else {
		log.Println("Discord Bot List Api is deactivated")
	}
}

func postStats(s *discordgo.Session, api *dblgo.Client) {
	log.Println("Posting stats to Discord Bot List")
	guildsNum := len(s.State.Guilds)
	log.Printf("Number of guilds %d", guildsNum)

	err := api.PostStats(guildsNum)
	if err != nil {
		log.Println("Failed to post stats to Discord Bot List", err)
	}
	log.Printf("Stats successfully posted")
}

func postStatsLoop(s *discordgo.Session, api dblgo.Client) {
	for {
		postStats(s, &api)
		time.Sleep(30 * time.Minute)
	}
}
