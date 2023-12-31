package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var (
	Token string
	Url   string
)

func init() {
	flag.StringVar(&Token, "t", "", "Discord Bot Token")
	flag.StringVar(&Url, "u", "", "IP API Url")
	flag.Parse()
}

func main() {
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	dg.Identify.Intents = discordgo.IntentsGuildMessages

	dg.AddHandler(messageCreate)

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
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
	resp, err := http.Get(Url)
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
