package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/kubastick/dblgo"
	"log"
	"time"
)

func postServerCountToDiscordBotApi(s *discordgo.Session) {
	if config.DiscordStatsSecret != "" && config.DiscordAppID != "" {
		dblApi := dblgo.NewDBLApi(config.DiscordStatsSecret)
		postStatsLoop(s, dblApi)
	} else {
		log.Println("Discord Bot List Api is deactivated")
	}
}

func postStats(s *discordgo.Session, api *dblgo.DBLApi) {
	log.Println("Posting stats to Discord Bot List")
	guildsNum := len(s.State.Guilds)
	log.Printf("Number of guilds %d", guildsNum)

	err := api.PostStatsSimple(guildsNum)
	if err != nil {
		log.Println("Failed to post stats to Discord Bot List", err)
	}
	log.Printf("Stats successfully posted")
}

func postStatsLoop(s *discordgo.Session, api dblgo.DBLApi) {
	for {
		postStats(s, &api)
		time.Sleep(30 * time.Minute)
	}
}
