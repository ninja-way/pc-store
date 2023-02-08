package transport

import (
	"github.com/gin-gonic/gin"
	"github.com/ninja-way/pc-store/internal/config"
	"github.com/ninja-way/pc-store/internal/models"
	"net/http"
)

// TODO: return info message if sign up not successful
func (h *Handler) signUp(c *gin.Context) {
	var signUp models.SignUp
	if err := c.BindJSON(&signUp); err != nil {
		config.LogDebug("signUp", err)
		c.Status(http.StatusBadRequest)
		return
	}

	if err := signUp.Validate(); err != nil {
		config.LogDebug("signUp", err)
		c.Status(http.StatusBadRequest)
		return
	}

	if err := h.userService.SignUp(c, signUp); err != nil {
		config.LogError("signUp", err)
		c.Status(http.StatusInternalServerError)
	}
}
