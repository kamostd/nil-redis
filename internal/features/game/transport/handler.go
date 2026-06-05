package transport

import (
	"context"
	"net/http"
	"nil-redis/internal/domain"
	"nil-redis/internal/features/game/transport/dto"
	"nil-redis/internal/features/game/transport/mapper"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type PlayerService interface {
	FetchPlayer(ctx context.Context, name string) (domain.Player, error)
	SavePlayer(ctx context.Context, name string) error
}

type PlayerHanders struct {
	service PlayerService
	logger  *zap.Logger
}

func New(service PlayerService, logger *zap.Logger) *PlayerHanders {
	return &PlayerHanders{
		service: service,
		logger:  logger,
	}
}

func (p PlayerHanders) GetPlayer(c *gin.Context) {
	ctx := c.Request.Context()

	name := c.Param("name")
	if strings.TrimSpace(name) == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "invalid path param",
		})

		return
	}

	player, err := p.service.FetchPlayer(ctx, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, mapper.PlayerToResponse(player))
}

func (p PlayerHanders) SavePlayer(c *gin.Context) {
	ctx := c.Request.Context()

	playerReq := dto.SavePlayerRequest{}

	if err := c.BindJSON(&playerReq); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "invalid body data",
		})

		return
	}

	p.service.SavePlayer(ctx)
}
