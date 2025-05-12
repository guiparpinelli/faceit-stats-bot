package main

import (
	"context"
	"database/sql"
	"log"

	"github.com/guiparpinelli/faceit-stats-bot/internal/infrastructure/db/sqlite"
	_ "github.com/mattn/go-sqlite3"
)

var ddl string

func run() error {
	ctx := context.Background()

	// Use a file-based SQLite database for persistence
	db, err := sql.Open("sqlite3", "file:./app.db")
	if err != nil {
		return err
	}
	defer db.Close()

	// Create tables
	if _, err := db.ExecContext(ctx, ddl); err != nil {
		return err
	}

	queries := sqlite.New(db)

	// Retrieve all players
	players, err := queries.FindAllPlayers(ctx)
	if err != nil {
		return err
	}
	log.Println(players)
	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
