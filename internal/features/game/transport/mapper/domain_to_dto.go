package mapper

import (
	"nil-redis/internal/domain"
	"nil-redis/internal/features/game/transport/dto"
)

func PlayerToResponse(player domain.Player) dto.GetPlayerResponse {
	return dto.GetPlayerResponse{
		Name:          player.Name(),
		BestScore:     player.BestScore(),
		TotalGames:    player.TotalGames(),
		RegisteredAtq: player.RegisteredAt(),
	}
}
