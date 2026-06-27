package usecase

import (
	"context"
	"nil-redis/internal/core/apperrors"
	"nil-redis/internal/domain"
)

type PlayerGeter interface {
	GetPlayerByName(ctx context.Context, name string) (domain.Player, error)

	IsExists(ctx context.Context, name string) (bool, error)
}

type RankScoreGeter interface {
	GetTotalScoreByName(ctx context.Context, name string) (int, error)
	RankByName(ctx context.Context, name string) (int, error)
}

type ProfileService struct {
	player    PlayerGeter
	rankScore RankScoreGeter
}

func NewProfileService(player PlayerGeter, rank RankScoreGeter) *ProfileService {
	return &ProfileService{
		player:    player,
		rankScore: rank,
	}
}

func (p ProfileService) CollectProfile(ctx context.Context, name string) (domain.Profile, error) {
	isExists, err := p.player.IsExists(ctx, name)
	if err != nil {
		return domain.Profile{}, err
	}

	if !isExists {
		return domain.Profile{}, apperrors.NotExists
	}

	player, err := p.player.GetPlayerByName(ctx, name)
	if err != nil {
		return domain.Profile{}, err
	}

	totalScore, err := p.rankScore.GetTotalScoreByName(ctx, name)
	if err != nil {
		return domain.Profile{}, err
	}

	rank, err := p.rankScore.RankByName(ctx, name)
	if err != nil {
		return domain.Profile{}, err
	}

	return domain.NewProfile(
		name,
		rank,
		totalScore,
		player.TotalGames(),
		player.BestScore(),
		player.RegisteredAt(),
	), nil
}
