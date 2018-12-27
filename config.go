package main

import (
	"github.com/BurntSushi/toml"
)

type Config struct {
	DiscordSecret string `toml:"discord_secret"`
}

func loadConfig() Config {
	config := Config{}
	_, err := toml.DecodeFile("config.toml", &config)
	if err != nil {
		panic(err)
	}
	return config
}
