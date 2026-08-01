package mapper

import (
	"nil-redis/internal/domain"
	"nil-redis/internal/features/game/storage/model"
	"time"
)

func ModelPlayerToDomain(model model.Player, name string) (domain.Player, error) {
	registredAt, err := time.Parse(time.DateOnly, model.RegisteredAt)
	if err != nil {
		return domain.Player{}, err
	}

	return domain.New(name, model.TotalGames, model.BestScore, registredAt), nil
}
