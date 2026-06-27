package usecase

import (
	"context"
	"nil-redis/internal/domain"
	"time"
)

type GameRepo interface {
	IsExists(ctx context.Context, name string) (bool, error)
	UpdateBestScore(ctx context.Context, name string, score int) error
	UpdateGamesPlayed(ctx context.Context, name string) error
	WriteRegistrationDate(ctx context.Context, name string, registredAt time.Time) error
	GetBestScoreByName(ctx context.Context, name string) (int, error)
}

type LeaderboardSaverGeter interface {
	SaveOrUpdateMember(ctx context.Context, name string, score int) (int, error)
	RankByName(ctx context.Context, name string) (int, error)
}

type HistorySaver interface {
	SaveGame(ctx context.Context, historyNote string) error
}

type StatsUpdater interface {
	UpdateTotalGame(ctx context.Context) error
}

type GameService struct {
	player      GameRepo
	leaderboard LeaderboardSaverGeter
	history     HistorySaver
	stats       StatsUpdater
}

func NewGameService(
	player GameRepo,
	leaderboard LeaderboardSaverGeter,
	history HistorySaver,
	stats StatsUpdater,
) *GameService {
	return &GameService{
		player:      player,
		leaderboard: leaderboard,
		history:     history,
		stats:       stats,
	}
}

// TODO: добавить транзакцию
func (p GameService) Result(ctx context.Context, name string, score int) (domain.LeaderboardPlayer, error) {
	isExist, err := p.player.IsExists(ctx, name)
	if err != nil {
		return domain.LeaderboardPlayer{}, err
	}

	if !isExist {
		if err := p.player.WriteRegistrationDate(ctx, name, time.Now()); err != nil {
			return domain.LeaderboardPlayer{}, err
		}
	}

	currentBestScore, err := p.player.GetBestScoreByName(ctx, name)
	if err != nil {
		return domain.LeaderboardPlayer{}, err
	}

	if currentBestScore < score {
		if err := p.player.UpdateBestScore(ctx, name, score); err != nil {
			return domain.LeaderboardPlayer{}, err
		}
	}

	if err := p.player.UpdateGamesPlayed(ctx, name); err != nil {
		return domain.LeaderboardPlayer{}, err
	}

	note := domain.NewNote(name, score, time.Now())

	if err := p.history.SaveGame(ctx, note.CreateNote()); err != nil {
		return domain.LeaderboardPlayer{}, err
	}

	if err := p.stats.UpdateTotalGame(ctx); err != nil {
		return domain.LeaderboardPlayer{}, err
	}

	totalScore, err := p.leaderboard.SaveOrUpdateMember(ctx, name, score)
	if err != nil {
		return domain.LeaderboardPlayer{}, err
	}

	rank, err := p.leaderboard.RankByName(ctx, name)
	if err != nil {
		return domain.LeaderboardPlayer{}, err
	}

	leaderboardNote := domain.NewLeaderboard(name, rank, totalScore)

	return leaderboardNote, nil
}
