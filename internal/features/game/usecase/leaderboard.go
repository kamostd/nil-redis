package usecase

import (
	"context"
	"nil-redis/internal/domain"
)

type LeaderboardLister interface {
	ListTopPlayerWithScore(ctx context.Context, top int) ([]domain.LeaderboardPlayer, error)
}

type LeaderboardService struct {
	repo LeaderboardLister
}

func NewLeaderboardService(repo LeaderboardLister) *LeaderboardService {
	return &LeaderboardService{
		repo: repo,
	}
}

func (l LeaderboardService) LeaderboardWithScore(ctx context.Context, top int) ([]domain.LeaderboardPlayer, error) {
	return l.repo.ListTopPlayerWithScore(ctx, top-1)
}
