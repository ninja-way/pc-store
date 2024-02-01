package transport

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ninja-way/pc-store/internal/config"
	"github.com/ninja-way/pc-store/internal/models"
	"net/http"
)

// TODO: return info message if sign up not successful

// @Summary		Sign Up
// @Description	Sign up new user
// @Tags			auth
// @Accept			json
// @Produce		json
// @Param			request	body	models.SignUp	true	"sign up info"
// @Success		200		"OK"
// @Failure		400		"bad request body"
// @Failure		500		"sign up error"
// @Router			/auth/sign-up [post]
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

// @Summary		Sign In
// @Description	Sign in user
// @Tags			auth
// @Accept			json
// @Produce		json
// @Param			request	body	models.SignIn	true	"sign in info"
// @Success		200		"OK"
// @Failure		400		"bad request body"
// @Failure		500		"sign in error"
// @Router			/auth/sign-in [post]
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

	accessToken, refreshToken, err := h.userService.SignIn(c, signIn)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}

		config.LogError("signIn", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Header("Set-Cookie", fmt.Sprintf("refresh-token=%s; HttpOnly", refreshToken))
	c.JSON(http.StatusOK, map[string]string{"token": accessToken})
}

// @Summary		Refresh Tokens
// @Description	Refresh access and refresh tokens
// @Tags			auth
// @Accept			json
// @Produce		json
// @Success		200	"OK"
// @Failure		400	"bad request body"
// @Failure		500	"refresh tokens error"
// @Router			/auth/refresh [get]
func (h *Handler) refresh(c *gin.Context) {
	cookie, err := c.Cookie("refresh-token")
	if err != nil {
		config.LogDebug("refresh", err)
		c.Status(http.StatusBadRequest)
		return
	}

	accessToken, refreshToken, err := h.userService.RefreshTokens(c, cookie)
	if err != nil {
		config.LogDebug("refresh", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Header("Set-Cookie", fmt.Sprintf("refresh-token=%s; HttpOnly", refreshToken))
	c.JSON(http.StatusOK, map[string]string{"token": accessToken})
}
