package command

import (
	"fmt"
	"pi-hosting/internal/config"
	"pi-hosting/internal/util"

	"github.com/bwmarrin/discordgo"
)

var commands = []*discordgo.ApplicationCommand{
	{
		Name:        "status",
		Description: "List status of all registered services",
	},
	{
		Name:        "start",
		Description: "Command for starting a service",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "service-name",
				Description: "Service name",
				Required:    true,
			},
		},
	},
	{
		Name:        "stop",
		Description: "Command for stopping a service",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "service-name",
				Description: "Service name",
				Required:    true,
			},
		},
	},
}

var handlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate, c *config.Config){
	"status": statusCommandHandler,
	"start":  startCommandHandler,
	"stop":   stopCommandHandler,
}

func RegisterHandler(s *discordgo.Session, c *config.Config) {
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := handlers[i.ApplicationCommandData().Name]; ok {
			h(s, i, c)
		}
	})
}

func RegisterCommands(s *discordgo.Session) ([]*discordgo.ApplicationCommand, error) {
	createdCommands, err := s.ApplicationCommandBulkOverwrite(s.State.User.ID, "", commands)
	if err != nil {
		return nil, err
	}

	return createdCommands, nil
}

func statusCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate, c *config.Config) {
	msg := ""
	for _, service := range c.Services {
		status, err := util.CheckServiceStatus(service)
		if err != nil {
			status = "Error"
		}
		msg += fmt.Sprintf(
			"Service: %s\nStatus: %s\n",
			service,
			status,
		)
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
		},
	})
}

func startCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate, c *config.Config) {
	options := i.ApplicationCommandData().Options

	service := options[0].StringValue()

	found := false
	msg := "Service not found"

	for s := range c.Services {
		if c.Services[s] == service {
			found = true
		}
	}

	if found {
		_, err := util.StartService(service)
		if err != nil {
			msg = fmt.Sprintf(
				"Error starting service: %s",
				service,
			)
		} else {
			msg = fmt.Sprintf(
				"Starting service: %s",
				service,
			)
		}
	} else {
		msg = "Service not found"
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
		},
	})
}

func stopCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate, c *config.Config) {
	options := i.ApplicationCommandData().Options

	service := options[0].StringValue()

	found := false
	msg := "Service not found"

	for s := range c.Services {
		if c.Services[s] == service {
			found = true
		}
	}

	if found {
		_, err := util.StopService(service)
		if err != nil {
			msg = fmt.Sprintf(
				"Error stopping service: %s",
				service,
			)
		} else {
			msg = fmt.Sprintf(
				"Stopping service: %s",
				service,
			)
		}
	} else {
		msg = "Service not found"
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
		},
	})
}
