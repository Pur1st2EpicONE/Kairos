// Package handler provides HTTP handlers for the Kairos application.
// It sets up routes, static file serving, and HTML templates for the web interface.
package handler

import (
	"Kairos/internal/service"
	"errors"
	"net/http"
	"strings"

	v1 "Kairos/internal/handler/v1"

	"github.com/gin-gonic/gin"
	"github.com/wb-go/wbf/ginext"
)

const templatePath = "web/templates/index.html"

func NewHandler(service *service.Service) http.Handler {

	handler := ginext.New("")

	handler.Use(ginext.Recovery())

	apiV1 := handler.Group("/api/v1")
	handlerV1 := v1.NewHandler(*service)

	auth := apiV1.Group("/auth")
	auth.POST("/sign-up", handlerV1.SignUp)
	auth.POST("/sign-in", handlerV1.SignIn)

	protected := apiV1.Group("/")
	protected.Use(JWT(service.AuthService))

	protected.POST("/events", handlerV1.CreateEvent)

	protected.GET("/notify", handlerV1.GetNotification)
	protected.DELETE("/notify", handlerV1.CancelNotification)

	return handler

}

func JWT(service service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			v1.RespondError(c, errors.New("empty auth header"))
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 {
			v1.RespondError(c, errors.New("invalid auth header"))
			return
		}

		userID, err := service.ParseToken(parts[1])
		if err != nil {
			v1.RespondError(c, errors.New("invalid token"))
			return
		}

		c.Set("user_id", userID)
		c.Next()

	}
}
