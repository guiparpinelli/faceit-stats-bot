package faceit

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type CS2Stats struct {
	SkillLevel int `json:"skill_level"`
	FaceitElo  int `json:"faceit_elo"`
}

type Games struct {
	CS2 CS2Stats `json:"cs2"`
}

type Player struct {
	Id          string    `json:"player_id"`
	Nickname    string    `json:"nickname"`
	Avatar      string    `json:"avatar"`
	Country     string    `json:"country"`
	Games       Games     `json:"games"`
	ActivatedAt time.Time `json:"activated_at"`
}

type PlayerRepository interface {
	Create(player *Player) (*Player, error)
	FindAll() ([]*Player, error)
	FindById(id uuid.UUID) (*Player, error)
	FindByNickname(nickname string) (*Player, error)
}

type PlayerService struct {
	client Client
	repo   PlayerRepository
}

func (p *PlayerService) GetTrackedPlayers() ([]*Player, error) {
	return p.repo.FindAll()
}

func (p *PlayerService) TrackPlayer(nickname string) (*Player, error) {
	// First check if player is already being tracked
	if player, _ := p.repo.FindByNickname(nickname); player != nil {
		return player, nil
	}

	// Try to fetch player from FaceitAPI
	requestUrl := fmt.Sprintf(p.client.baseUrl+"/players?nickname=%s", nickname)
	req, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		return nil, err
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var player Player
	if err := json.NewDecoder(resp.Body).Decode(&player); err != nil {
		return nil, err
	}

	// Add new tracked player
	if _, err := p.repo.Create(&player); err != nil {
		return nil, err
	}

	return &player, nil
}

/*
func (p *PlayerService) UntrackPlayer(nickname string) (*Player, error) {
	player, err := p.repo.FindByNickname(nickname)
	if err != nil {
		return nil, err
	}
	if player == nil {
		return nil, fmt.Errorf("Player %s not found", nickname)
	}

	if err := p.repo.Remove(player.Id); err != nil {
		return nil, err
	}

	return player, nil
}
*/
