package transport

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/ninja-way/pc-store/internal/models"
	"github.com/ninja-way/pc-store/internal/transport/middleware"

	_ "github.com/ninja-way/pc-store/docs"
	swaggerFiles "github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"

	"strconv"
)

/******** Transport layer *********/

// ComputersStore is service layer entity
type ComputersStore interface {
	GetComputers(context.Context) ([]models.PC, error)
	GetComputerByID(context.Context, int) (models.PC, error)
	AddComputer(context.Context, models.PC) (int, error)
	UpdateComputer(context.Context, int, models.PC) error
	DeleteComputer(context.Context, int) error
}

type Handler struct {
	service ComputersStore
}

func NewHandler(service ComputersStore) *Handler {
	return &Handler{
		service: service,
	}
}

// InitRouter setup endpoints
func (h *Handler) InitRouter() *gin.Engine {
	// disable debug info
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(middleware.Logging())

	r.GET("/computers", h.getComputers)

	comp := r.Group("/computer")
	{
		comp.PUT("", h.addComputer)
		comp.GET("/:id", h.getComputer)
		comp.POST("/:id", h.updateComputer)
		comp.DELETE("/:id", h.deleteComputer)
	}

	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return r
}

func parseID(s string) (int, error) {
	id, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}

	if id == 0 {
		return 0, errors.New("id can't be 0")
	}

	return id, nil
}
