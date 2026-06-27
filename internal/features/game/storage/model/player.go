package model

type Player struct {
	TotalGames   int    `redis:"games_played"`
	BestScore    int    `redis:"best_score"`
	RegisteredAt string `redis:"registered_at"`
}

func NewPlayer(registeredAt string, bestScore, gamesPlayed int) Player {
	return Player{
		TotalGames:   gamesPlayed,
		BestScore:    bestScore,
		RegisteredAt: registeredAt,
	}
}
