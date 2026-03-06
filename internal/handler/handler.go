package handler

import (
	"Kairos/internal/config"
	"Kairos/internal/errs"
	"Kairos/internal/service"
	"context"
	"html/template"
	"net/http"
	"strings"

	v1 "Kairos/internal/handler/v1"

	"github.com/gin-gonic/gin"
	"github.com/wb-go/wbf/ginext"
)

const templatePath = "web/templates/index.html"
const AuthHeader = "Authorization"

func NewHandler(config config.Server, service *service.Service) http.Handler {

	handler := ginext.New("")

	handler.Use(ginext.Recovery())

	apiV1 := handler.Group("/api/v1")
	handlerV1 := v1.NewHandler(config, *service)

	auth := apiV1.Group("/auth")
	auth.POST("/sign-up", handlerV1.SignUp)
	auth.POST("/sign-in", handlerV1.SignIn)

	apiV1.GET("/events/:id", handlerV1.GetInfo)
	apiV1.GET("/", homePage(template.Must(template.ParseFiles(templatePath)), service.CoreService))

	protected := apiV1.Group("/")
	protected.Use(authJWT(service.AuthService))

	protected.POST("/events", handlerV1.CreateEvent)
	protected.POST("/events/:id/book", handlerV1.CreateBooking)
	protected.POST("/events/:id/confirm", handlerV1.ConfirmBooking)

	return handler

}

func authJWT(service service.AuthService) gin.HandlerFunc {

	return func(c *gin.Context) {

		authHeader := c.GetHeader(AuthHeader)
		if authHeader == "" {
			v1.RespondError(c, errs.ErrEmptyAuthHeader)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 {
			v1.RespondError(c, errs.ErrInvalidAuthHeader)
			return
		}

		userID, err := service.ParseToken(parts[1])
		if err != nil {
			v1.RespondError(c, err)
			return
		}

		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), "userID", userID))
		c.Next()

	}

}

func homePage(tmpl *template.Template, service service.CoreService) func(c *ginext.Context) {
	return func(c *ginext.Context) {
		events := service.GetAllEvents(c.Request.Context())
		c.Header("Content-Type", "text/html")
		if err := tmpl.Execute(c.Writer, map[string]any{"Events": events}); err != nil {
			c.String(http.StatusInternalServerError, errs.ErrInternal.Error())
		}
	}
}
