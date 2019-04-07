package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
)

type Config struct {
	DiscordSecret string `toml:"discord_secret"`
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
		DiscordSecret: "YOUR_DISCORD_BOT_TOKEN",
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
