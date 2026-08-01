package transport

import (
	"context"
	"net/http"
	"nil-redis/internal/domain"
	"nil-redis/internal/features/game/transport/dto"
	"nil-redis/internal/features/game/transport/mapper"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type GameService interface {
	Result(ctx context.Context, name string, score int) (domain.LeaderboardPlayer, error)
}

type HistoryServiece interface {
	ListHistoryNotes(ctx context.Context) ([]domain.HistoryNote, error)
}

type LeaderboardService interface {
	LeaderboardWithScore(ctx context.Context, top int) ([]domain.LeaderboardPlayer, error)
}

type ProfileCollector interface {
	CollectProfile(ctx context.Context, name string) (domain.Profile, error)
}

type StatsService interface {
	TotalGames(ctx context.Context) (int, error)
}

type HTTPHandler struct {
	game        GameService
	history     HistoryServiece
	leaderboard LeaderboardService
	profile     ProfileCollector
	stats       StatsService

	logger *zap.Logger
}

func NewHTTPHandler(
	game GameService,
	history HistoryServiece,
	leaderboard LeaderboardService,
	profile ProfileCollector,
	stats StatsService,
	logger *zap.Logger,
) *HTTPHandler {
	return &HTTPHandler{
		game:        game,
		history:     history,
		leaderboard: leaderboard,
		profile:     profile,
		stats:       stats,
		logger:      logger,
	}
}

func (h HTTPHandler) Routing(group *gin.RouterGroup) {
	group.GET("/leaderboard", h.getLeaders)
	group.GET("/player/:name", h.getProfile)
	group.GET("/history", h.getHistory)
	group.GET("/stats", h.getTotalGame)

	gameGroup := group.Group("/game")
	gameGroup.POST("/result", h.result)

}

// Result godoc
// @Summary      Save game result
// @Description  register new player and save score or save score in exist player
// @Tags         game
// @Accept       json
// @Produce      json
// @Param        game_result  body dto.ResultRequest  true  "Game result with score"
// @Success      200  {object}  dto.ResultResponse
// @Router       /game/result [post]“
func (h HTTPHandler) result(c *gin.Context) {
	ctx := c.Request.Context()

	req := dto.ResultRequest{}
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"messaage": "invalid request body"})

		return
	}

	leaderboard, err := h.game.Result(ctx, req.Name, req.Score)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"messaage": err.Error()})

		return
	}

	c.JSON(http.StatusOK, mapper.LeaderboartToResultResponse(leaderboard))
}

// GetHistory godoc
// @Summary      get games history
// @Description  shows the last 10 games played
// @Accept       json
// @Produce      json
// @Success      200  {object}  []dto.GetHistoryResponse
// @Router       /history  [get]
func (h HTTPHandler) getHistory(c *gin.Context) {
	ctx := c.Request.Context()

	notes, err := h.history.ListHistoryNotes(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"messaage": err.Error()})

		return
	}

	c.JSON(http.StatusOK, mapper.HistoryNotesToGetHistoryResponse(notes))
}

// GetLeaderboard godoc
// @Summary      get leaderboard
// @Description  shows the specified number of leaders
// @Accept       json
// @Produce      json
// @Param        top query int true "leaders count"
// @Success      200  {object}  []dto.GetLeadersResponse
// @Router       /leaderboard [get]
func (h HTTPHandler) getLeaders(c *gin.Context) {
	ctx := c.Request.Context()

	topS := c.Query("top")

	top, err := strconv.Atoi(topS)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"messaage": "invalid request query"})

		return
	}

	leaders, err := h.leaderboard.LeaderboardWithScore(ctx, top)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"messaage": err.Error()})

		return
	}

	c.JSON(http.StatusOK, mapper.LeaderboardsToGetLeadersResponse(leaders))
}

// GetProfile godoc
// @Summary      show palayer profile
// @Description  get player profile by name
// @Tags         player
// @Accept       json
// @Produce      json
// @Param        name   path      string  true  "Player name"
// @Success      200  {object}  dto.GetProfileResponse
// @Router       /player/{name} [get]
func (h HTTPHandler) getProfile(c *gin.Context) {
	ctx := c.Request.Context()

	name := c.Params.ByName("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "empty name"})

		return
	}

	profile, err := h.profile.CollectProfile(ctx, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"messaage": err.Error()})

		return
	}

	c.JSON(http.StatusOK, mapper.ProfileToGetProfileResponse(profile))
}

// GetTotalGame godoc
// @Summary      show total_game count
// @Description  get int
// @Accept       json
// @Produce      json
// @Success      200  {object}  dto.GetTotalGameResponse
// @Router       /stats [get]
func (h HTTPHandler) getTotalGame(c *gin.Context) {
	ctx := c.Request.Context()

	games, err := h.stats.TotalGames(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"messaage": err.Error()})

		return
	}

	c.JSON(http.StatusOK, dto.GetTotalGameResponse{
		TotalGames: games,
	})
}
