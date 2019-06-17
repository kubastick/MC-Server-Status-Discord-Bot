package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
)

type Config struct {
	DiscordSecret      string `toml:"discord_secret"`
	DiscordStatsSecret string `toml:"discord_stats_secret"`
	DiscordAppID       string `toml:"discord_app_id"`
	DisableStatus      bool   `toml:"disable_status"`
}

func loadConfig() Config {
	config := Config{}
	_, err := toml.DecodeFile("config.toml", &config)
	if err != nil {
		createConfig()
		fmt.Println("Created new config file - please fill in discord token")
		os.Exit(1)
	}
	return config
}

func createConfig() {
	fillInCfg := Config{
		DiscordSecret: "NTg5MDk2OTEwNzY2MzQyMTU3.XQeSUQ.ywUBaZYCMECDjLygSCpDgkYZh4M",
	}
	cfgFile, err := os.Create("config.toml")
	if err != nil {
		panic("Failed to create config.toml")
	}
	defer cfgFile.Close()

	encoder := toml.NewEncoder(cfgFile)
	err = encoder.Encode(fillInCfg)
	if err != nil {
		panic(err)
	}
}
