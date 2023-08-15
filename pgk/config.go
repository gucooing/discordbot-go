package pgk

import (
	"encoding/json"
	"errors"
	"os"
)

type Config struct {
	Addr            string `json:"Addr"`
	Port            string `json:"Port"`
	Path            string `json:"Path"`
	McHost          string `json:"McHost"`
	DiscordBotToken string `json:"DiscordBotToken"`
	GuildID         string `json:"GuildID"`
}

var CONF *Config = nil

func GetConfig() *Config {
	return CONF
}

var FileNotExist = errors.New("config file not found")

func LoadConfig() error {
	filePath := "./config.json"
	if len(os.Args) > 1 {
		filePath = os.Args[1]
	}
	f, err := os.Open(filePath)
	if err != nil {
		return FileNotExist
	}
	defer func() {
		_ = f.Close()
	}()
	c := new(Config)
	d := json.NewDecoder(f)
	if err := d.Decode(c); err != nil {
		return err
	}
	CONF = c
	return nil
}

var DefaultConfig = &Config{
	Addr:            "0.0.0.0",
	Port:            "8080",
	Path:            "discord",
	McHost:          "127.0.0.1:19132",
	DiscordBotToken: "1234567890",
	GuildID:         "",
}
