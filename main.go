package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"gopkg.in/yaml.v2"
)

func main() {
	config, err := loadConfig()
	if err != nil {
		log.Fatal("error loading config,", err)
	}

	dg, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		log.Fatal("error creating Discord session,", err)
	}

	dg.Identify.Intents = discordgo.IntentsGuildMessages

	dg.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		messageCreate(s, m, &config)
	})

	err = dg.Open()
	if err != nil {
		log.Fatal("error opening connection,", err)
	}

	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate, config *Config) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	mentioned := false
	for _, mention := range m.Mentions {
		if mention.ID == s.State.User.ID {
			mentioned = true
		}
	}

	if !mentioned {
		return
	}

	msg := "error getting IP address"
	resp, err := http.Get(config.URL)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, msg)
		return
	}
	defer resp.Body.Close()
	ip, err := io.ReadAll(resp.Body)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, msg)
		return
	}

	msg = fmt.Sprintf("IP address: %s", string(ip))

	s.ChannelMessageSend(m.ChannelID, msg)
}

type Config struct {
	Token    string   `yaml:"token"`
	URL      string   `yaml:"ip_url"`
	Services []string `yaml:"services,flow"`
}

func loadConfig() (Config, error) {
	configFile, err := os.ReadFile("config.yaml")
	if err != nil {
		return Config{}, err
	}

	config := Config{}
	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}
