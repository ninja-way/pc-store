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

type Users interface {
	SignUp(context.Context, models.SignUp) error
}

type Handler struct {
	compService ComputersStore
	userService Users
}

func NewHandler(compService ComputersStore, userService Users) *Handler {
	return &Handler{
		compService: compService,
		userService: userService,
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

	// auth
	auth := r.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
	}

	// computers
	r.GET("/computers", h.getComputer)

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
