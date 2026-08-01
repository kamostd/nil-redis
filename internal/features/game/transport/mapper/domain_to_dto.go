package mapper

import (
	"nil-redis/internal/domain"
	"nil-redis/internal/features/game/transport/dto"
	"time"
)

func LeaderboartToResultResponse(lb domain.LeaderboardPlayer) dto.ResultResponse {
	return dto.ResultResponse{
		Name:          lb.Name(),
		NewTotalScore: lb.Score(),
		Rank:          lb.Rank(),
	}
}

func LeaderboardsToGetLeadersResponse(leaders []domain.LeaderboardPlayer) []dto.GetLeadersResponse {
	dtos := make([]dto.GetLeadersResponse, 0, len(leaders))

	for _, leader := range leaders {
		dtos = append(dtos, dto.GetLeadersResponse{
			Rank:  leader.Rank(),
			Name:  leader.Name(),
			Score: leader.Score(),
		})
	}

	return dtos
}

func HistoryNotesToGetHistoryResponse(notes []domain.HistoryNote) []dto.GetHistoryResponse {
	dtos := make([]dto.GetHistoryResponse, 0, len(notes))

	for _, note := range notes {
		dtos = append(dtos, dto.GetHistoryResponse{
			Name:     note.Name(),
			Score:    note.Score(),
			PlayedAt: note.CreatedAt().Format("2006-01-02T15:04:05"),
		})
	}

	return dtos
}

func ProfileToGetProfileResponse(profile domain.Profile) dto.GetProfileResponse {
	return dto.GetProfileResponse{
		Name:         profile.Name(),
		Rank:         profile.Rank(),
		TotalScore:   profile.TotalScore(),
		GamesPlayed:  profile.GamePlayed(),
		BestScore:    profile.BestScore(),
		RegisteredAt: profile.RegisteredAt().Format(time.DateOnly),
	}
}
