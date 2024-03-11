package command

import (
	"fmt"
	"pi-hosting/internal/config"

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
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Hey there! Congratulations, you just executed your first slash command",
		},
	})
}

func startCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate, c *config.Config) {
	// Access options in the order provided by the user.
	options := i.ApplicationCommandData().Options

	// Or convert the slice into a map
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	// This example stores the provided arguments in an []interface{}
	// which will be used to format the bot's response
	margs := make([]interface{}, 0, len(options))
	msgformat := "You learned how to use command options! " +
		"Take a look at the value(s) you entered:\n"

	// Get the value from the option map.
	// When the option exists, ok = true
	if option, ok := optionMap["string-option"]; ok {
		// Option values must be type asserted from interface{}.
		// Discordgo provides utility functions to make this simple.
		margs = append(margs, option.StringValue())
		msgformat += "> string-option: %s\n"
	}
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		// Ignore type for now, they will be discussed in "responses"
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf(
				msgformat,
				margs...,
			),
		},
	})
}

func stopCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate, c *config.Config) {
	// Access options in the order provided by the user.
	options := i.ApplicationCommandData().Options

	// Or convert the slice into a map
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	// This example stores the provided arguments in an []interface{}
	// which will be used to format the bot's response
	margs := make([]interface{}, 0, len(options))
	msgformat := "You learned how to use command options! " +
		"Take a look at the value(s) you entered:\n"

	// Get the value from the option map.
	// When the option exists, ok = true
	if option, ok := optionMap["string-option"]; ok {
		// Option values must be type asserted from interface{}.
		// Discordgo provides utility functions to make this simple.
		margs = append(margs, option.StringValue())
		msgformat += "> string-option: %s\n"
	}
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		// Ignore type for now, they will be discussed in "responses"
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf(
				msgformat,
				margs...,
			),
		},
	})
}
