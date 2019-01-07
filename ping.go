package main

import (
	"github.com/bwmarrin/discordgo"
	"log"
)

func handlePing(s *discordgo.Session, m *discordgo.MessageCreate) {
	log.Printf("Received ping message \"%s\" from %s\n", m.Content, m.Author.Username)
	response := "Pong!"
	log.Printf("Responding with %s", response)
	_, err := s.ChannelMessageSend(m.ChannelID, response)
	if err != nil {
		log.Printf("Failed to send message: %s", err.Error())
	}
}
