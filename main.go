package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
)

// Bot represents our Discord bot instance
type Bot struct {
	session         *discordgo.Session
	commands        []*discordgo.ApplicationCommand
	commandHandlers map[string]CommandHandlerFunc
}

// CommandHandlerFunc defines the function signature for command handlers
type CommandHandlerFunc func(s *discordgo.Session, i *discordgo.InteractionCreate)

// Command definitions
var commands = []*discordgo.ApplicationCommand{
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

func main() {
	token := flag.String("t", "", "Discord Bot Token")
	flag.Parse()

	// Create a new bot instance
	bot, err := NewBot(*token)
	if err != nil {
		log.Fatalf("Error creating bot: %v", err)
	}
	defer bot.Close()

	// Register commands and start the bot
	if err := bot.RegisterCommands(); err != nil {
		log.Fatalf("Error registering commands: %v", err)
	}

	// Wait for interrupt signal to gracefully shutdown
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt)
	<-sc

	log.Println("Shutdown complete")
}

// NewBot creates a new bot instance
func NewBot(token string) (*Bot, error) {
	// Create Discord session
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, fmt.Errorf("error creating Discord session: %w", err)
	}

	bot := &Bot{
		session:         dg,
		commands:        commands,
		commandHandlers: make(map[string]CommandHandlerFunc),
	}

	// Register handlers
	bot.RegisterHandlers()

	// Open connection to Discord
	if err := dg.Open(); err != nil {
		return nil, fmt.Errorf("error opening connection: %w", err)
	}

	// Add ready handler
	dg.AddHandler(bot.handleReady)

	// Add interaction handler
	dg.AddHandler(bot.handleSlashCommand)

	return bot, nil
}

// RegisterHandlers registers all command handlers
func (b *Bot) RegisterHandlers() {
	b.commandHandlers["ping"] = b.handlePing
	b.commandHandlers["help"] = b.handleHelp
	b.commandHandlers["track"] = b.handleTrack
	b.commandHandlers["untrack"] = b.handleUntrack
	b.commandHandlers["list"] = b.handleList
}

// RegisterCommands registers all slash commands with Discord
func (b *Bot) RegisterCommands() error {
	for _, command := range b.commands {
		_, err := b.session.ApplicationCommandCreate(b.session.State.User.ID, "", command)
		if err != nil {
			return fmt.Errorf("error creating slash command %q: %w", command.Name, err)
		}
		log.Printf("Registered command: %s", command.Name)
	}
	return nil
}

// Close closes the Discord session
func (b *Bot) Close() {
	if b.session != nil {
		b.session.Close()
	}
}

// Event handlers
func (b *Bot) handleReady(s *discordgo.Session, event *discordgo.Ready) {
	log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
}

func (b *Bot) handleSlashCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	commandName := i.ApplicationCommandData().Name
	if handler, exists := b.commandHandlers[commandName]; exists {
		handler(s, i)
	}
}

// Command handlers
func (b *Bot) handlePing(s *discordgo.Session, i *discordgo.InteractionCreate) {
	latency := s.HeartbeatLatency().Milliseconds()

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Pong! %dms", latency),
		},
	})
}

func (b *Bot) handleHelp(s *discordgo.Session, i *discordgo.InteractionCreate) {
	embed := &discordgo.MessageEmbed{
		Title:       "Help",
		Description: "Here are all the available commands:",
		Color:       0x00ff00,
		Fields:      make([]*discordgo.MessageEmbedField, 0, len(b.commands)),
	}

	for _, cmd := range b.commands {
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

func (b *Bot) handleTrack(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return
}

func (b *Bot) handleUntrack(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return
}

func (b *Bot) handleList(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return
}
