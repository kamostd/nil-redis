package domain

import "time"

type Player struct {
	name         string
	bestScore    int
	totalGame    int
	registeredAt time.Time
}

func New(name string, totalGame, bestScore int, registeredAt time.Time) Player {
	return Player{
		name:         name,
		bestScore:    bestScore,
		registeredAt: registeredAt,
		totalGame:    totalGame,
	}
}

func (p Player) Name() string {
	return p.name
}

func (p Player) BestScore() int {
	return p.bestScore
}
func (p Player) TotalGames() int {
	return p.totalGame
}
func (p Player) RegisteredAt() time.Time {
	return p.registeredAt
}
