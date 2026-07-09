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

type TxManager interface {
	Do(ctx context.Context, tx func(context.Context) error) error
}

type GameService struct {
	player      GameRepo
	leaderboard LeaderboardSaverGeter
	history     HistorySaver
	stats       StatsUpdater
	trm         TxManager
}

func NewGameService(
	player GameRepo,
	leaderboard LeaderboardSaverGeter,
	history HistorySaver,
	stats StatsUpdater,
	trm TxManager,
) *GameService {
	return &GameService{
		player:      player,
		leaderboard: leaderboard,
		history:     history,
		stats:       stats,
		trm:         trm,
	}
}

// TODO: добавить транзакцию
func (p GameService) Result(ctx context.Context, name string, score int) (domain.LeaderboardPlayer, error) {
	var leaderboard_note domain.LeaderboardPlayer

	if err := p.trm.Do(ctx, func(ctx context.Context) error {
		isExist, err := p.player.IsExists(ctx, name)
		if err != nil {
			return err
		}

		if !isExist {
			if err := p.player.WriteRegistrationDate(ctx, name, time.Now()); err != nil {
				return err
			}
		}

		currentBestScore, err := p.player.GetBestScoreByName(ctx, name)
		if err != nil {
			return err
		}

		if currentBestScore < score {
			if err := p.player.UpdateBestScore(ctx, name, score); err != nil {
				return err
			}
		}

		if err := p.player.UpdateGamesPlayed(ctx, name); err != nil {
			return err
		}

		note := domain.NewNote(name, score, time.Now())

		if err := p.history.SaveGame(ctx, note.CreateNote()); err != nil {
			return err
		}

		if err := p.stats.UpdateTotalGame(ctx); err != nil {
			return err
		}

		totalScore, err := p.leaderboard.SaveOrUpdateMember(ctx, name, score)
		if err != nil {
			return err
		}

		rank, err := p.leaderboard.RankByName(ctx, name)
		if err != nil {
			return err
		}

		leaderboard_note = domain.NewLeaderboard(name, rank, totalScore)

		return nil
	}); err != nil {
		return domain.LeaderboardPlayer{}, err
	}

	return leaderboard_note, nil
}
