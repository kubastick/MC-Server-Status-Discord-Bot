package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

var (
	config = loadConfig()
)

func main() {
	configureLogger()

	session := connectToDiscord(config.DiscordSecret)
	session.AddHandler(messageRouter)
	go postServerCountToDiscordBotApi(session)
	go statusLoop(session)
	defer session.Close()

	log.Println("Minecraft status bot is ready!")
	// Wait for something that look like sigkill
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	log.Println("Received SIGKILL, exiting")
}

func connectToDiscord(secret string) *discordgo.Session {
	discord, err := discordgo.New("Bot " + secret)
	if err != nil {
		log.Fatalf("Failed to connect, reasen: %s", err.Error())
	}
	err = discord.Open()
	if err != nil {
		log.Fatalf("Fatal to establish websocket connection with discord.io, reason: %s", err.Error())
	}
	return discord
}

func messageRouter(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Do not parse own and bot messages
	if s.State.User.ID == m.Author.ID || m.Author.Bot {
		return
	}

	// Get content of message
	userMessage := m.Content

	// Commands routing
	if strings.HasPrefix(userMessage, "!status ") {
		handleStatus(s, m)
	}

	if strings.HasPrefix(userMessage, "!ping") {
		handlePing(s, m)
	}

	if strings.HasPrefix(userMessage, "!help") {
		handleHelp(s, m)
	}
}

func configureLogger() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func statusLoop(s *discordgo.Session) {
	if config.DisableStatus {
		log.Println("Status has been disabled in config file")
		return
	}
	for {
		log.Println("Updating status")
		statusData := discordgo.UpdateStatusData{
			Status: "online",
			Game: &discordgo.Game{
				Type: discordgo.GameTypeWatching,
				Name: fmt.Sprintf("%d servers", len(s.State.Guilds)),
			},
		}
		err := s.UpdateStatusComplex(statusData)
		if err != nil {
			log.Println("Warning, failed to update discord status")
		}
		time.Sleep(5 * time.Minute)
	}
}
