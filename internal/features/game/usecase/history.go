package usecase

import (
	"context"
	"nil-redis/internal/domain"
)

type HistoryLister interface {
	ListHistory(ctx context.Context) ([]domain.HistoryNote, error)
}

type HistoryService struct {
	repo HistoryLister
}

func NewHistoryService(repo HistoryLister) *HistoryService {
	return &HistoryService{
		repo: repo,
	}
}

func (h HistoryService) ListHistoryNotes(ctx context.Context) ([]domain.HistoryNote, error) {
	return h.repo.ListHistory(ctx)
}
