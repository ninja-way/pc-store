package transport

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/ninja-way/pc-store/internal/config"
	"github.com/ninja-way/pc-store/internal/models"
	"github.com/ninja-way/pc-store/internal/transport/middleware"

	_ "github.com/ninja-way/pc-store/docs"
	swaggerFiles "github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
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
	compService ComputersStore
}

func NewHandler(compService ComputersStore) *Handler {
	return &Handler{
		compService: compService,
	}
}

// InitRouter setup endpoints
func (h *Handler) InitRouter(cfg *config.Config) *gin.Engine {
	// disable debug info
	if cfg.Environment == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(middleware.Logging())

	r.GET("/computers", h.getComputers)

	comp := r.Group("/computer")
	{
		comp.POST("", h.addComputer)
		comp.GET("/:id", h.getComputer)
		comp.PUT("/:id", h.updateComputer)
		comp.DELETE("/:id", h.deleteComputer)
	}

	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return r
}
