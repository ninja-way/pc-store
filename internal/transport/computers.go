package transport

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/ninja-way/pc-store/internal/models"
	"github.com/ninja-way/pc-store/internal/service"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type IDResponse struct {
	ID int `json:"id" example:"1"`
}

//	@Summary		Get Computers
//	@Description	Get all pc from database
//	@Tags			computers
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	[]models.PC
//	@Failure		500	"get pcs from database error"
//	@Router			/computers [get]
func (h *Handler) getComputers(c *gin.Context) {
	comps, err := h.service.GetComputers(c)
	if err != nil {
		log.WithFields(log.Fields{
			"handler": "getComputers",
			"problem": "get all pc from service",
		}).Error(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, comps)
}

//	@Summary		Add Computer
//	@Description	Add new pc from request body to database
//	@Tags			computer
//	@Accept			json
//	@Produce		json
//	@Param			request	body		models.PC	true	"computer and its accessories"
//	@Success		201		{object}	IDResponse
//	@Failure		400		"bad request body"
//	@Failure		500		"add pc to database error"
//	@Router			/computer [post]
func (h *Handler) addComputer(c *gin.Context) {
	var newPC models.PC
	if err := c.BindJSON(&newPC); err != nil {
		log.WithFields(log.Fields{
			"handler": "addComputer",
			"problem": "read or unmarshal request body",
		}).Debug(err)
		c.Status(http.StatusBadRequest)
		return
	}

	id, err := h.service.AddComputer(c, newPC)
	if err != nil {
		if errors.Is(err, service.ErrPriceTooHigh) || errors.Is(err, service.ErrFewComponents) {
			log.WithFields(log.Fields{
				"handler": "addComputer",
				"problem": "bad request body",
			}).Debug(err)
			c.Status(http.StatusBadRequest)
			return
		}

		log.WithFields(log.Fields{
			"handler": "addComputer",
			"problem": "add pc to service",
		}).Error(err)
		c.Status(http.StatusInternalServerError)
	}

	c.JSON(http.StatusCreated, IDResponse{id})
}

//	@Summary		Get Computer
//	@Description	Get pc from database by id
//	@Tags			computer
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Computer ID"
//	@Success		200	{object}	models.PC
//	@Failure		400	"bad id passed"
//	@Failure		404	"pc with passed id not found"
//	@Router			/computer/{id} [get]
func (h *Handler) getComputer(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		log.WithFields(log.Fields{
			"handler": "getComputer",
			"problem": "bad id passed",
		}).Debug(err)
		c.Status(http.StatusBadRequest)
		return
	}

	pc, err := h.service.GetComputerByID(c, id)
	if err != nil {
		log.WithFields(log.Fields{
			"handler": "getComputer",
			"problem": "pc with passed id not found",
		}).Debug(err)
		c.Status(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, pc)
}

//	@Summary		Update Computer
//	@Description	Update existing pc in database by id
//	@Tags			computer
//	@Accept			json
//	@Produce		json
//	@Param			id		path	int			true	"Computer ID"
//	@Param			request	body	models.PC	true	"new computer or some new accessories"
//	@Success		200		"pc updated"
//	@Failure		400		"bad id passed"
//	@Failure		400		"bad request body"
//	@Failure		400		"pc with passed id not found"
//	@Router			/computer/{id} [put]
func (h *Handler) updateComputer(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		log.WithFields(log.Fields{
			"handler": "updateComputer",
			"problem": "bad id passed",
		}).Debug(err)
		c.Status(http.StatusBadRequest)
		return
	}

	var newPC models.PC
	if err = c.BindJSON(&newPC); err != nil {
		log.WithFields(log.Fields{
			"handler": "updateComputer",
			"problem": "read or unmarshal request body",
		}).Debug(err)
		c.Status(http.StatusBadRequest)
		return
	}

	if err = h.service.UpdateComputer(c, id, newPC); err != nil {
		if errors.Is(err, service.ErrPriceTooHigh) {
			log.WithFields(log.Fields{
				"handler": "updateComputer",
				"problem": "specified price too high",
			}).Debug(err)
			c.Status(http.StatusBadRequest)
			return
		}
		log.WithFields(log.Fields{
			"handler": "updateComputer",
			"problem": "pc with passed id not found",
		}).Debug(err)
		c.Status(http.StatusNotFound)
	}
}

//	@Summary		Delete Computer
//	@Description	Delete pc from database by id
//	@Tags			computer
//	@Accept			json
//	@Produce		json
//	@Param			id	path	int	true	"Computer ID"
//	@Success		200	"pc deleted"
//	@Failure		400	"bad id passed"
//	@Failure		400	"pc with passed id not found"
//	@Router			/computer/{id} [delete]
func (h *Handler) deleteComputer(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		log.WithFields(log.Fields{
			"handler": "deleteComputer",
			"problem": "bad id passed",
		}).Debug(err)
		c.Status(http.StatusBadRequest)
		return
	}

	if err = h.service.DeleteComputer(c, id); err != nil {
		log.WithFields(log.Fields{
			"handler": "deleteComputer",
			"problem": "pc with passed id not found",
		}).Debug(err)
		c.Status(http.StatusNotFound)
	}
}
