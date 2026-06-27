package domain

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

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

type HistoryNote struct {
	name      string
	score     int
	createdAt time.Time
}

func NewNote(name string, score int, endTime time.Time) HistoryNote {
	return HistoryNote{
		name:      name,
		score:     score,
		createdAt: endTime,
	}
}

func NewNoteByString(note string) (HistoryNote, error) {
	noteUnits := strings.Split(note, ":")

	createdAt, err := time.Parse("2006-01-02T15:04:05", strings.Join(noteUnits[2:], ":"))
	if err != nil {
		return HistoryNote{}, err
	}

	score, err := strconv.Atoi(noteUnits[1])
	if err != nil {
		return HistoryNote{}, err
	}

	return HistoryNote{
		name:      noteUnits[0],
		score:     score,
		createdAt: createdAt,
	}, nil

}

func (h HistoryNote) Name() string {
	return h.name
}

func (h HistoryNote) Score() int {
	return h.score
}

func (h HistoryNote) CreatedAt() time.Time {
	return h.createdAt
}

func (h HistoryNote) CreateNote() string {
	return fmt.Sprintf(
		"%s:%d:%s",
		h.name,
		h.score,
		h.createdAt.Format("2006-01-02T15:04:05"),
	)
}

type LeaderboardPlayer struct {
	name  string
	rank  int
	score int
}

func NewLeaderboard(name string, rank, score int) LeaderboardPlayer {
	return LeaderboardPlayer{
		name:  name,
		rank:  rank,
		score: score,
	}
}

func (l LeaderboardPlayer) Name() string {
	return l.name
}

func (l LeaderboardPlayer) Rank() int {
	return l.rank
}

func (l LeaderboardPlayer) Score() int {
	return l.score
}

type Profile struct {
	name        string
	rank        int
	totalScore  int
	gamePlayed  int
	bestScore   int
	registredAt time.Time
}

func NewProfile(
	name string,
	rank, totalScore, gamePlayed, bestScore int,
	registredAt time.Time,
) Profile {
	return Profile{
		name:        name,
		rank:        rank,
		totalScore:  totalScore,
		gamePlayed:  gamePlayed,
		bestScore:   bestScore,
		registredAt: registredAt,
	}
}

func (p Profile) Name() string {
	return p.name
}

func (p Profile) Rank() int {
	return p.rank
}

func (p Profile) TotalScore() int {
	return p.totalScore
}
func (p Profile) GamePlayed() int {
	return p.gamePlayed
}

func (p Profile) BestScore() int {
	return p.bestScore
}

func (p Profile) RegisteredAt() time.Time {
	return p.registredAt
}
