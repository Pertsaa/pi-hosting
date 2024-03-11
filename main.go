package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"pi-hosting/internal/command"
	"pi-hosting/internal/config"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

func main() {
	config, err := config.Load()
	if err != nil {
		log.Fatal("error loading config,", err)
	}

	s, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		log.Fatal("error creating Discord session,", err)
	}

	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	command.RegisterHandler(s, &config)

	err = s.Open()
	if err != nil {
		log.Fatal("error opening connection,", err)
	}
	defer s.Close()

	_, err = command.RegisterCommands(s)
	if err != nil {
		log.Fatal("error registering commands,", err)
	}

	fmt.Println("Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	log.Println("Gracefully shutting down.")
}
