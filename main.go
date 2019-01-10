package main

import (
	"MinecraftServerStatusBot/mcsrvstat"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io"
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

func checkMinecraftServer(address string) (image io.Reader, statusText string, err error) {
	status, err := mcsrvstat.Query(address)
	if err != nil {
		return nil, "", err
	}

	// Try generate image
	img, err := status.GenerateStatusImage()
	if err != nil {
		log.Println("Failed to generate image: " + err.Error())
	}

	// And generate text
	result := fmt.Sprintf("Players online:  %d \\ %d\n", status.Players.Online, status.Players.Max)
	if len(status.Motd.Clean) > 0 {
		result += fmt.Sprintf("MOTD: %s\n", status.Motd.Clean[0])
	}
	result += fmt.Sprintf("Version: %s\n", status.Version)
	playerCount := len(status.Players.List)
	if playerCount > 0 && playerCount < 10 {
		result += "Player list:\n"
		for _, p := range status.Players.List {
			result += p + "\n"
		}
	}
	return &img, result, nil
}
