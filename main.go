package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	token         = flag.String("t", "", "Discord Bot Token")
	faceitApiKey  = flag.String("k", "", "FACEIT Api Key")
	faceitBaseURL = "https://open.faceit.com/data/v4"
	commands      = []*discordgo.ApplicationCommand{
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
			Description: "Track a FACEIT player by their nickname. e.g.: `/track <nickname>`",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "player",
					Description: "The FACEIT nickname",
					Required:    true,
				},
			},
		},
		{
			Name:        "untrack",
			Description: "Untrack a FACEIT player by their nickname. e.g.: `/untrack <nickname>`",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "player",
					Description: "The FACEIT nickname",
					Required:    true,
				},
			},
		},
		{
			Name:        "list",
			Description: "List all players that are being tracked.",
		},
	}

	// In-memory storage for tracked players
	trackedPlayers = struct {
		sync.RWMutex
		players map[string]Player
	}{
		players: make(map[string]Player),
	}
)

// addTrackedPlayer adds a player to the tracked players map
func addTrackedPlayer(player Player) {
	trackedPlayers.Lock()
	defer trackedPlayers.Unlock()
	trackedPlayers.players[strings.ToLower(player.Nickname)] = player
}

// removeTrackedPlayer removes a player from the tracked players map
func removeTrackedPlayer(nickname string) bool {
	trackedPlayers.Lock()
	defer trackedPlayers.Unlock()
	nickname = strings.ToLower(nickname)
	if _, exists := trackedPlayers.players[nickname]; exists {
		delete(trackedPlayers.players, nickname)
		return true
	}
	return false
}

// getTrackedPlayers returns a copy of all tracked players
func getTrackedPlayers() []Player {
	trackedPlayers.RLock()
	defer trackedPlayers.RUnlock()
	players := make([]Player, 0, len(trackedPlayers.players))
	for _, player := range trackedPlayers.players {
		players = append(players, player)
	}
	return players
}

// isPlayerTracked checks if a player is already being tracked
func isPlayerTracked(nickname string) bool {
	trackedPlayers.RLock()
	defer trackedPlayers.RUnlock()
	_, exists := trackedPlayers.players[strings.ToLower(nickname)]
	return exists
}

// populateDevData adds some sample players to the in-memory storage for development
func populateDevData() {
	samplePlayers := []Player{
		{
			Nickname: "s1mple",
			Avatar:   "https://assets.faceit-cdn.net/avatars/b5b0fe3f-25d9-4845-a3e9-6ba95866734e_1550799664341.jpg",
			Country:  "UA",
			Games: Games{
				CS2: CS2Stats{
					SkillLevel: 10,
					FaceitElo:  3555,
				},
			},
			ActivatedAt: time.Date(2014, 4, 22, 0, 0, 0, 0, time.UTC),
		},
		{
			Nickname: "ZywOo",
			Avatar:   "https://assets.faceit-cdn.net/avatars/2c998a4c-2228-4e63-b42a-b3ed60c9972d_1551292270057.jpg",
			Country:  "FR",
			Games: Games{
				CS2: CS2Stats{
					SkillLevel: 10,
					FaceitElo:  3498,
				},
			},
			ActivatedAt: time.Date(2015, 11, 8, 0, 0, 0, 0, time.UTC),
		},
		{
			Nickname: "NiKo",
			Avatar:   "https://assets.faceit-cdn.net/avatars/9a4e0b94-9454-4f73-9446-7f0e639a5911_1550799665045.jpg",
			Country:  "BA",
			Games: Games{
				CS2: CS2Stats{
					SkillLevel: 10,
					FaceitElo:  3402,
				},
			},
			ActivatedAt: time.Date(2014, 8, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			Nickname: "ropz",
			Avatar:   "https://assets.faceit-cdn.net/avatars/9a4e0b94-9454-4f73-9446-7f0e639a5911_1550799665045.jpg",
			Country:  "EE",
			Games: Games{
				CS2: CS2Stats{
					SkillLevel: 10,
					FaceitElo:  3350,
				},
			},
			ActivatedAt: time.Date(2016, 1, 30, 0, 0, 0, 0, time.UTC),
		},
		{
			Nickname: "m0NESY",
			Avatar:   "https://assets.faceit-cdn.net/avatars/9a4e0b94-9454-4f73-9446-7f0e639a5911_1550799665045.jpg",
			Country:  "RU",
			Games: Games{
				CS2: CS2Stats{
					SkillLevel: 10,
					FaceitElo:  3289,
				},
			},
			ActivatedAt: time.Date(2020, 5, 12, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, player := range samplePlayers {
		addTrackedPlayer(player)
	}
	log.Printf("Populated %d sample players for development", len(samplePlayers))
}

func main() {
	flag.Parse()

	// Populate development data
	populateDevData()

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

type CS2Stats struct {
	SkillLevel int `json:"skill_level"`
	FaceitElo  int `json:"faceit_elo"`
}

type Games struct {
	CS2 CS2Stats `json:"cs2"`
}

type Player struct {
	Nickname    string    `json:"nickname"`
	Avatar      string    `json:"avatar"`
	Country     string    `json:"country"`
	Games       Games     `json:"games"`
	ActivatedAt time.Time `json:"activated_at"`
}

func handleTrack(s *discordgo.Session, i *discordgo.InteractionCreate) {
	nickname := i.ApplicationCommandData().Options[0].StringValue()
	msg := ""

	// Check if player is already being tracked
	if isPlayerTracked(nickname) {
		// Get the tracked player from storage
		trackedPlayers.RLock()
		player := trackedPlayers.players[strings.ToLower(nickname)]
		trackedPlayers.RUnlock()

		embed := &discordgo.MessageEmbed{
			Title: fmt.Sprintf("%s :flag_%s:", player.Nickname, strings.ToLower(player.Country)),
			Color: 0x00ff00, // Green color
			Description: fmt.Sprintf("Level: %d\nELO: %d\nMember since: %s",
				player.Games.CS2.SkillLevel,
				player.Games.CS2.FaceitElo,
				player.ActivatedAt.Format("January 2, 2006")),
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: player.Avatar,
			},
		}
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{embed},
			},
		})
		return
	}

	// Log the request details
	requestURL := fmt.Sprintf(faceitBaseURL+"/players?nickname=%s", nickname)

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		msg = fmt.Sprintf("Error creating request: %v", err)
		log.Printf("Request creation error: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+*faceitApiKey)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		msg = fmt.Sprintf("Error sending request: %v", err)
		log.Printf("Request execution error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		msg = fmt.Sprintf("HTTP Error: %s", resp.Status)
		log.Printf("Response body: %s", string(body))
	}

	var player Player
	if err := json.NewDecoder(resp.Body).Decode(&player); err != nil {
		msg = fmt.Sprintf("error decoding response: %v", err)
		log.Printf("Decoding error: %v", err)
	}

	if msg != "" {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: msg,
			},
		})
		return
	}

	// Add player to tracked players
	addTrackedPlayer(player)

	embed := &discordgo.MessageEmbed{
		Title: fmt.Sprintf("%s :flag_%s:", player.Nickname, strings.ToLower(player.Country)),
		Color: 0x00ff00, // Green color
		Description: fmt.Sprintf("Level: %d\nELO: %d\nMember since: %s",
			player.Games.CS2.SkillLevel,
			player.Games.CS2.FaceitElo,
			player.ActivatedAt.Format("January 2, 2006")),
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: player.Avatar,
		},
	}
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})
}

func handleUntrack(s *discordgo.Session, i *discordgo.InteractionCreate) {
	nickname := i.ApplicationCommandData().Options[0].StringValue()
	if removed := removeTrackedPlayer(nickname); removed {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("Successfully untracked player %s!", nickname),
			},
		})
	} else {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("Player %s was not being tracked!", nickname),
			},
		})
	}
}

func handleList(s *discordgo.Session, i *discordgo.InteractionCreate) {
	players := getTrackedPlayers()
	if len(players) == 0 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "No players are currently being tracked!",
			},
		})
		return
	}

	// Sort players by ELO
	sort.Slice(players, func(i, j int) bool {
		return players[i].Games.CS2.FaceitElo > players[j].Games.CS2.FaceitElo
	})

	embed := &discordgo.MessageEmbed{
		Title:       "Tracked Players",
		Description: fmt.Sprintf("Currently tracking %d players", len(players)),
		Color:       0x00ff00, // Green color
	}

	var listContent strings.Builder
	listContent.WriteString("# - Player - ELO\n")
	listContent.WriteString("----------------\n")
	for i, player := range players {
		listContent.WriteString(fmt.Sprintf("%d. %s - %d\n",
			i+1,
			player.Nickname,
			player.Games.CS2.FaceitElo))
	}

	embed.Description = fmt.Sprintf("%s\n```%s```", embed.Description, listContent.String())

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})
}

var commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	"ping":    handlePing,
	"help":    handleHelp,
	"track":   handleTrack,
	"untrack": handleUntrack,
	"list":    handleList,
}
