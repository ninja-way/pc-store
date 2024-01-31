package transport

import (
	"errors"
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

func (h *Handler) signIn(c *gin.Context) {
	var signIn models.SignIn
	if err := c.BindJSON(&signIn); err != nil {
		config.LogDebug("signIn", err)
		c.Status(http.StatusBadRequest)
		return
	}

	if err := signIn.Validate(); err != nil {
		config.LogDebug("signIn", err)
		c.Status(http.StatusBadRequest)
		return
	}

	token, err := h.userService.SignIn(c, signIn)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}

		config.LogError("signIn", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, map[string]string{"token": token})
}
