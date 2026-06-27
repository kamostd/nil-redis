package usecase

import "context"

type StatsGeter interface {
	GetTotalGames(ctx context.Context) (int, error)
}

type StatsService struct {
	repo StatsGeter
}

func NewStatsService(repo StatsGeter) *StatsService {
	return &StatsService{
		repo: repo,
	}
}

func (s StatsService) TotalGames(ctx context.Context) (int, error) {
	return s.repo.GetTotalGames(ctx)
}
