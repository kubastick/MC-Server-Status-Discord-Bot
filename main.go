package main

import (
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func main() {
	config := loadConfig()

	session := connectToDiscord(config.DiscordSecret)
	session.AddHandler(messageRouter)
	postServerCountToDiscordBotApi(session)
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
		os.Exit(-1)
	}
	err = discord.Open()
	if err != nil {
		log.Fatalf("Fatal to establish websocket connection with discord.io, reason: %s", err.Error())
		os.Exit(-1)
	}
	return discord
}

func messageRouter(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Do not parse own messages
	if s.State.User.ID == m.Author.ID {
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
