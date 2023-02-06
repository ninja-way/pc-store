package transport

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/ninja-way/pc-store/internal/models"
	"github.com/ninja-way/pc-store/internal/transport/middleware"
	"log"
	"net/http"
	"strconv"
)

type ComputersStore interface {
	GetComputers(context.Context) ([]models.PC, error)
	GetComputerByID(context.Context, int) (models.PC, error)
	AddComputer(context.Context, models.PC) error
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

// InitRouter setup handlers
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

	return r
}

// Get method which returns all pc from database
func (h *Handler) getComputers(c *gin.Context) {
	comps, err := h.service.GetComputers(c)
	if err != nil {
		log.Println("GetComputers error:", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, comps)
}

// Put method which add new pc from request body to database
func (h *Handler) addComputer(c *gin.Context) {
	var newPC models.PC
	if err := c.BindJSON(&newPC); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	if err := h.service.AddComputer(c, newPC); err != nil {
		log.Println("AddComputer error:", err)
		c.Status(http.StatusInternalServerError)
	}
}

// Get method which return pc from database by id
func (h *Handler) getComputer(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	pc, err := h.service.GetComputerByID(c, id)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, pc)
}

// Post method which update existing pc in database by id and new data from request body
func (h *Handler) updateComputer(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	var newPC models.PC
	if err = c.BindJSON(&newPC); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	if err = h.service.UpdateComputer(c, id, newPC); err != nil {
		c.Status(http.StatusNotFound)
	}
}

// Delete method which delete pc from database by id
func (h *Handler) deleteComputer(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	if err = h.service.DeleteComputer(c, id); err != nil {
		c.Status(http.StatusNotFound)
	}
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
