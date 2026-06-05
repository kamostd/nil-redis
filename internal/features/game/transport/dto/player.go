package dto

import "time"

type ErrorResponse struct {
	Message string `json:"message"`
}

type GetPlayerResponse struct {
	Name          string    `json:"name"`
	BestScore     int       `json:"best_score"`
	TotalGames    int       `json:"total_games"`
	RegisteredAtq time.Time `json:"registered_atq"`
}

type SavePlayerRequest struct {
	Name       string `json:"name"`
	BestScore  int    `json:"best_score"`
	TotalGames int    `json:"total_games"`
}
