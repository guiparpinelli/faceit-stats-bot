package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
)

var (
	token    = flag.String("t", "", "Discord Bot Token")
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "ping",
			Description: "Shows the bot's latency",
		},
		{
			Name:        "help",
			Description: "Shows the bot's commands",
		},
		{
			Name:        "track",
			Description: "Track a FACEIT player. e.g.: `/track <nickname> OR <profile URL>`",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "player",
					Description: "The FACEIT nickname or profile URL",
					Required:    true,
				},
			},
		},
		{
			Name:        "untrack",
			Description: "Untrack a FACEIT player. e.g.: `/untrack <nickname> OR <profile URL>`",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "player",
					Description: "The FACEIT nickname or profile URL",
					Required:    true,
				},
			},
		},
		{
			Name:        "list",
			Description: "List all players that are being tracked.",
		},
	}
)

func main() {
	flag.Parse()

	// Create a new Discord session
	dg, err := discordgo.New("Bot " + *token)
	if err != nil {
		log.Fatalf("Error creating Discord session: %v", err)
	}

	// Defer closing the Discord session
	defer dg.Close()

	// Register the messageCreate func as a callback for MessageCreate events
	dg.AddHandler(ready)

	// Open a websocket connection to Discord
	err = dg.Open()
	if err != nil {
		log.Fatalf("Error opening connection: %v", err)
	}

	// Register all slash commands
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, command := range commands {
		cmd, err := dg.ApplicationCommandCreate(dg.State.User.ID, "", command)
		if err != nil {
			log.Fatalf("Error creating slash command %q: %v", command.Name, err)
		}
		registeredCommands[i] = cmd
		log.Printf("Registered command: %s", command.Name)
	}

	// Add handler for slash commands
	dg.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		handleSlashCommand(s, i)
	})

	// Wait here until CTRL-C or other term signal is received
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt)
	<-sc

	log.Println("Shutdown complete")
}

func ready(s *discordgo.Session, event *discordgo.Ready) {
	log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
}

func handleSlashCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	if handler, exists := commandHandlers[i.ApplicationCommandData().Name]; exists {
		handler(s, i)
	}
}

func handlePing(s *discordgo.Session, i *discordgo.InteractionCreate) {
	latency := s.HeartbeatLatency().Milliseconds()

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Pong! %dms", latency),
		},
	})
}

func handleHelp(s *discordgo.Session, i *discordgo.InteractionCreate) {
	embed := &discordgo.MessageEmbed{
		Title:       "Help",
		Description: "Here are all the available commands:",
		Color:       0x00ff00,
		Fields:      make([]*discordgo.MessageEmbedField, 0, len(commands)),
	}

	for _, cmd := range commands {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "/" + cmd.Name,
			Value:  cmd.Description,
			Inline: false,
		})
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})
}

var commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	"ping": handlePing,
	"help": handleHelp,
	"track": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintln("Not implemented"),
			},
		})
	},
	"untrack": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintln("Not implemented"),
			},
		})
	},
	"list": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintln("Not implemented"),
			},
		})
	},
}
