package main

import (
	"github.com/bwmarrin/discordgo"
	"log"
)

func handleHelp(s *discordgo.Session, m *discordgo.MessageCreate) {
	const helpMessage = `:white_check_mark: List of commands:


:bulb: !status <server_address> [--text] - Shows graphical server status [--text - Text instead of graphics]
:bulb: !ping - "Pong"

:hammer: Examples:
!status hypixel.net --text
!ping
!status mistylands.net`

	log.Printf("User \"%s\" asked for help using \"%s\" command\n", m.Author.Username, m.Content)
	// Reply with helpMessage
	_, err := s.ChannelMessageSend(m.ChannelID, helpMessage)
	if err != nil {
		log.Println("Failed to send help message: " + err.Error())
	} else {
		log.Println("Successfully send help message")
	}
}
