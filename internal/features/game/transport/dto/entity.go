package dto

type ErrorResponse struct {
	Message string `json:"message"`
}

type ResultRequest struct {
	Name  string `json:"name"`
	Score int    `json:"score"`
}

type ResultResponse struct {
	Name          string `json:"name"`
	NewTotalScore int    `json:"new_total_score"`
	Rank          int    `json:"rank"`
}

type GetLeadersResponse struct {
	Rank  int    `json:"rank"`
	Name  string `json:"name"`
	Score int    `json:"score"`
}

type GetHistoryResponse struct {
	Name     string `json:"name"`
	Score    int    `json:"score"`
	PlayedAt string `json:"played_at"`
}

type GetProfileResponse struct {
	Name         string `json:"name"`
	Rank         int    `json:"rank"`
	TotalScore   int    `json:"total_score"`
	GamesPlayed  int    `json:"games_played"`
	BestScore    int    `json:"best_score"`
	RegisteredAt string `json:"registered_at"`
}

type GetTotalGameResponse struct {
	TotalGames int `json:"total_games"`
}
