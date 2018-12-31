package main

import (
	"MinecraftServerStatusBot/mcsrvstat"
	"fmt"
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
	defer session.Close()

	log.Println("Minecraft status bot is ready!")
	// Wait for something that look like sigkill
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
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
	userMessage := m.Content
	// Check if message is command
	if strings.HasPrefix(userMessage, "!status ") {
		log.Printf("Parsing mesage: %s\n", userMessage)
		_, err := s.ChannelMessageSend(m.ChannelID, "Ok, I'm going to check this minecraft server IP!")
		if err != nil {
			log.Fatalf("Failed to send message %s\n", err.Error())
		}

		serverIP := strings.Replace(userMessage, "!status ", "", -1)
		log.Printf("Checking %s\n", serverIP)

		result, err := checkMinecraftServer(serverIP)
		if err != nil {
			responseMessage := fmt.Sprintf("Sorry, but I can't find minecraft server with these ip :c")
			log.Printf("Failed to server check IP: %s", err.Error())
			_, err = s.ChannelMessageSend(m.ChannelID, responseMessage)
			if err != nil {
				log.Fatalf("Failed to send message %s\n", err.Error())
			}
			return
		}

		_, err = s.ChannelMessageSend(m.ChannelID, "Here we go!")
		if err != nil {
			log.Fatalf("Failed to send message %s\n", err.Error())
		}

		_, err = s.ChannelMessageSend(m.ChannelID, result)
		if err != nil {
			log.Fatalf("Failed to send message %s\n", err.Error())
		}

		log.Printf("Everything went ok (%s)", serverIP)
	}
}

func checkMinecraftServer(address string) (string, error) {
	status, err := mcsrvstat.Query(address)
	if err != nil {
		return "", err
	}

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
	return result, nil
}
