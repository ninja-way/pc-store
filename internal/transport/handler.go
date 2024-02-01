package transport

import (
	"context"
	"github.com/gin-gonic/gin"
	_ "github.com/ninja-way/pc-store/docs"
	"github.com/ninja-way/pc-store/internal/config"
	"github.com/ninja-way/pc-store/internal/models"
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
	SignIn(context.Context, models.SignIn) (string, string, error)
	ParseToken(context.Context, string) (int64, error)
	RefreshTokens(context.Context, string) (string, string, error)
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
	r.Use(Logging())

	// auth
	auth := r.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.GET("/sign-in", h.signIn)
		auth.GET("/refresh", h.refresh)
	}

	// computers
	comp := r.Group("/computers")
	comp.Use(Auth())
	{
		comp.GET("", h.getComputers)
		comp.POST("", h.addComputer)
		comp.GET("/:id", h.getComputer)
		comp.PUT("/:id", h.updateComputer)
		comp.DELETE("/:id", h.deleteComputer)
	}

	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return r
}
