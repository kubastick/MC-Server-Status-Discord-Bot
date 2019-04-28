package main

import (
	"MinecraftServerStatusBot/mcsrvstat"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io"
	"log"
	"strings"
)

func handleStatus(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Get content out of message
	userMessage := m.Content

	log.Printf("User \"%s\" asked for server status using \"%s\" command", m.Author.Username, userMessage)
	_, err := s.ChannelMessageSend(m.ChannelID, "Ok, I'm going to check this minecraft server IP!")
	if err != nil {
		log.Printf("Failed to send message %s\n", err.Error())
	}
	// Show, that we are working
	_ = s.ChannelTyping(m.ChannelID)

	// Get server IP
	serverIP := strings.Replace(userMessage, "!status ", "", -1)
	serverIP = strings.Replace(serverIP, "--text", "", -1)
	serverIP = strings.TrimSpace(serverIP)

	log.Printf("Checking %s\n", serverIP)

	imageResult, result, err := checkMinecraftServer(serverIP)
	if err != nil {
		responseMessage := fmt.Sprintf("Sorry, but I can't find minecraft server with these ip :c")
		log.Printf("Failed to server check IP: %s\n", err.Error())
		_, err = s.ChannelMessageSend(m.ChannelID, responseMessage)
		if err != nil {
			// Failed to even send checking message
			log.Printf("Failed to send message %s\n", err.Error())
			return
		}
		return
	}

	_, err = s.ChannelMessageSend(m.ChannelID, "Here we go!")
	if err != nil {
		log.Printf("Failed to send message %s\n", err.Error())
	}

	// If user not decided otherwise, try sending image first
	if !strings.Contains(userMessage, "--text") {
		// Yeah user not decided otherwise, we are trying to send image
		if imageResult != nil {
			_, err := s.ChannelFileSend(m.ChannelID, "result.png", imageResult)
			if err != nil {
				log.Println("Failed to send image, using text as fallback")
			} else {
				// Successfully send image, so do not send text
				log.Println("Successfully sended message")
				return
			}
		}
	}

	// Otherwise send text message
	_, err = s.ChannelMessageSend(m.ChannelID, result)
	if err != nil {
		log.Printf("Failed to send message %s\n", err.Error())
		return
	}

	log.Printf("Everything went ok (%s)\n", serverIP)
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
