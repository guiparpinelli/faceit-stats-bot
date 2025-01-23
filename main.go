package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
)

func main() {
	// Parse command line flags
	var token string
	flag.StringVar(&token, "t", "", "Discord Bot Token")
	flag.Parse()

	if token == "" {
		log.Fatal("No token provided. Use -t flag to provide your Discord bot token")
	}

	// Create a new Discord session
	dg, err := discordgo.New("Bot " + token)
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

	// Register slash command
	command := &discordgo.ApplicationCommand{
		Name:        "ping",
		Description: "Shows the bot's latency",
	}

	_, err = dg.ApplicationCommandCreate(dg.State.User.ID, "", command)
	if err != nil {
		log.Fatalf("Error creating slash command: %v", err)
	}

	// Add handler for slash commands
	dg.AddHandler(handleSlashCommand)

	// Wait here until CTRL-C or other term signal is received
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt)
	<-sc
}

func ready(s *discordgo.Session, event *discordgo.Ready) {
	log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
}

func handleSlashCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	switch i.ApplicationCommandData().Name {
	case "ping":
		latency := s.HeartbeatLatency().Milliseconds()

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("Pong! Latency: %dms", latency),
			},
		})
	}
}
