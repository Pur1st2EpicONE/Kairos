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

const indexPath = "web/templates/index.html"
const addPath = "web/templates/create_event.html"
const loginPath = "web/templates/login.html"
const signupPath = "web/templates/signup.html"
const eventPath = "web/templates/event.html"

const header = "Authorization"

func NewHandler(config config.Server, service *service.Service) http.Handler {

	handler := ginext.New("")

	handler.Use(ginext.Recovery())

	handler.Static("/static", "./web/static")

	handler.GET("/", homePage(template.Must(template.ParseFiles(indexPath)), service.CoreService))
	handler.GET("/add", renderPage(template.Must(template.ParseFiles(addPath))))
	handler.GET("/login", renderPage(template.Must(template.ParseFiles(loginPath))))
	handler.GET("/signup", renderPage(template.Must(template.ParseFiles(signupPath))))
	handler.GET("/events/:id", eventPage(template.Must(template.ParseFiles(eventPath)), service.CoreService))

	apiV1 := handler.Group("/api/v1")
	handlerV1 := v1.NewHandler(config, *service)

	auth := apiV1.Group("/auth")
	auth.POST("/sign-up", handlerV1.SignUp)
	auth.POST("/sign-in", handlerV1.SignIn)

	apiV1.GET("/events/:id", handlerV1.GetInfo)

	protected := apiV1.Group("/")
	protected.Use(authJWT(service.AuthService))

	protected.POST("/events", handlerV1.CreateEvent)
	protected.POST("/events/:id/book", handlerV1.CreateBooking)
	protected.POST("/events/:id/confirm", handlerV1.ConfirmBooking)

	return handler

}

func authJWT(service service.AuthService) gin.HandlerFunc {

	return func(c *ginext.Context) {

		authHeader := c.GetHeader(header)
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

func renderPage(tmpl *template.Template) gin.HandlerFunc {
	return func(c *ginext.Context) {
		c.Header("Content-Type", "text/html")
		if err := tmpl.Execute(c.Writer, nil); err != nil {
			c.String(http.StatusInternalServerError, errs.ErrInternal.Error())
		}
	}
}

func homePage(tmpl *template.Template, service service.CoreService) gin.HandlerFunc {

	return func(c *ginext.Context) {

		events := service.GetAllEvents(c.Request.Context())
		eventsDTO := make([]v1.InfoResponseDTO, len(events))

		for i, e := range events {
			eventsDTO[i] = v1.InfoResponseDTO{
				ID:          e.ID,
				Title:       e.Title,
				Description: e.Description,
				Date:        e.Date,
				Seats:       e.Seats,
			}
		}

		c.Header("Content-Type", "text/html")
		if err := tmpl.Execute(c.Writer, map[string]any{"Events": eventsDTO}); err != nil {
			c.String(http.StatusInternalServerError, errs.ErrInternal.Error())
		}

	}

}

func eventPage(tmpl *template.Template, service service.CoreService) gin.HandlerFunc {

	return func(c *ginext.Context) {

		id := c.Param("id")
		event, err := service.GetInfo(c.Request.Context(), id)
		if err != nil {
			v1.RespondError(c, err)
			return
		}

		c.Header("Content-Type", "text/html")
		if err := tmpl.Execute(c.Writer, v1.InfoResponseDTO{
			ID:          id,
			Title:       event.Title,
			Description: event.Description,
			Date:        event.Date,
			Seats:       event.Seats,
		}); err != nil {
			c.String(http.StatusInternalServerError, errs.ErrInternal.Error())
		}

	}

}
