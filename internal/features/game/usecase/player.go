package usecase

import (
	"context"
	"nil-redis/internal/domain"
	"nil-redis/internal/features/game/storage/model"
	"time"
)

type PlayerRepo interface {
	SavePlayer(ctx context.Context, name string, bestScore int, totalGames int, registredAt time.Time) error
	FindPlayerByName(ctx context.Context, name string) (model.Player, error)
}

type PlayerService struct {
	repo PlayerRepo
}

func New(repo PlayerRepo) *PlayerService {
	return &PlayerService{
		repo: repo,
	}
}

func (p PlayerService) SavePlayer(ctx context.Context, name string) error {
	return p.repo.SavePlayer(ctx, name, 0, 0, time.Now())
}

func (p PlayerService) FetchPlayer(ctx context.Context, name string) (domain.Player, error) {
	model, err := p.repo.FindPlayerByName(ctx, name)
	if err != nil {
		return domain.Player{}, err
	}

	date, err := time.Parse(time.DateOnly, model.RegisteredAt)
	if err != nil {
		return domain.Player{}, err
	}

	return domain.New(
		name,
		model.TotalGames,
		model.BestScore,
		date,
	), nil
}
