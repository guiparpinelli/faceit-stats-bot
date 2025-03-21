package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
)

func main() {
	// Parse command line flags
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
