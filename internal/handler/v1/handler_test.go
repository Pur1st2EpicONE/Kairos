package v1

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"Kairos/internal/config"
	"Kairos/internal/errs"
	"Kairos/internal/models"
	"Kairos/internal/service"
	"Kairos/internal/service/mocks"

	"github.com/stretchr/testify/require"
	"github.com/wb-go/wbf/ginext"
	"go.uber.org/mock/gomock"
)

const (
	testUUID    = "123e4567-e89b-12d3-a456-426614174000"
	invalidUUID = "lqhfjlqhwfljkhqwklf"
	testUserID  = int64(1)
	testEventID = testUUID
)

func TestHandler_SignUp(t *testing.T) {

	controller := gomock.NewController(t)
	defer controller.Finish()

	mockAuth := mocks.NewMockAuthService(controller)
	mockCore := mocks.NewMockCoreService(controller)
	mockService := service.Service{AuthService: mockAuth, CoreService: mockCore}

	header := NewHandler(config.Server{}, mockService)

	router := ginext.New("")
	router.POST("/sign-up", header.SignUp)

	t.Run("invalid JSON binding", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/sign-up", bytes.NewBufferString(`{"login":"test"`))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		require.Equal(t, http.StatusBadRequest, resp.Code)
		require.Contains(t, resp.Body.String(), errs.ErrInvalidJSON.Error())
	})

	t.Run("service CreateUser error", func(t *testing.T) {
		mockAuth.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(int64(0), errors.New("internal error"))
		req := httptest.NewRequest(http.MethodPost, "/sign-up", bytes.NewBufferString(`{"login":"test","password":"pass123"}`))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		require.Equal(t, http.StatusInternalServerError, resp.Code)
	})

	t.Run("user already exists", func(t *testing.T) {
		mockAuth.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(int64(0), errs.ErrUserAlreadyExists)
		req := httptest.NewRequest(http.MethodPost, "/sign-up", bytes.NewBufferString(`{"login":"test","password":"pass123"}`))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		require.Equal(t, http.StatusConflict, resp.Code)
		require.Contains(t, resp.Body.String(), errs.ErrUserAlreadyExists.Error())
	})

	t.Run("CreateToken error", func(t *testing.T) {
		mockAuth.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(testUserID, nil)
		mockAuth.EXPECT().CreateToken(testUserID).Return("", errors.New("token generation failed"))
		req := httptest.NewRequest(http.MethodPost, "/sign-up", bytes.NewBufferString(`{"login":"test","password":"pass123"}`))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		require.Equal(t, http.StatusInternalServerError, resp.Code)
	})

	t.Run("success", func(t *testing.T) {
		mockAuth.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(testUserID, nil)
		mockAuth.EXPECT().CreateToken(testUserID).Return("token123", nil)
		req := httptest.NewRequest(http.MethodPost, "/sign-up", bytes.NewBufferString(`{"login":"test","password":"pass123"}`))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		require.Equal(t, http.StatusOK, resp.Code)
		require.Contains(t, resp.Body.String(), "token123")
	})

}

func TestHandler_SignIn(t *testing.T) {

	controller := gomock.NewController(t)
	defer controller.Finish()

	mockAuth := mocks.NewMockAuthService(controller)
	mockCore := mocks.NewMockCoreService(controller)
	mockService := service.Service{AuthService: mockAuth, CoreService: mockCore}

	header := NewHandler(config.Server{}, mockService)

	router := ginext.New("")
	router.POST("/sign-in", header.SignIn)

	t.Run("invalid JSON binding", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/sign-in", bytes.NewBufferString(`{"login":"test"`))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		require.Equal(t, http.StatusBadRequest, resp.Code)
		require.Contains(t, resp.Body.String(), errs.ErrInvalidJSON.Error())
	})

	t.Run("service GetUserId error", func(t *testing.T) {
		mockAuth.EXPECT().GetUserId(gomock.Any(), gomock.Any()).Return(int64(0), errs.ErrInvalidCredentials)
		req := httptest.NewRequest(http.MethodPost, "/sign-in", bytes.NewBufferString(`{"login":"test","password":"wrong"}`))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		require.Equal(t, http.StatusUnauthorized, resp.Code)
		require.Contains(t, resp.Body.String(), errs.ErrInvalidCredentials.Error())
	})

	t.Run("CreateToken error", func(t *testing.T) {
		mockAuth.EXPECT().GetUserId(gomock.Any(), gomock.Any()).Return(testUserID, nil)
		mockAuth.EXPECT().CreateToken(testUserID).Return("", errors.New("token generation failed"))
		req := httptest.NewRequest(http.MethodPost, "/sign-in", bytes.NewBufferString(`{"login":"test","password":"pass123"}`))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		require.Equal(t, http.StatusInternalServerError, resp.Code)
	})

	t.Run("success", func(t *testing.T) {
		mockAuth.EXPECT().GetUserId(gomock.Any(), gomock.Any()).Return(testUserID, nil)
		mockAuth.EXPECT().CreateToken(testUserID).Return("token123", nil)
		req := httptest.NewRequest(http.MethodPost, "/sign-in", bytes.NewBufferString(`{"login":"test","password":"pass123"}`))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		require.Equal(t, http.StatusOK, resp.Code)
		require.Contains(t, resp.Body.String(), "token123")
	})

}

func TestHandler_GetInfo(t *testing.T) {

	controller := gomock.NewController(t)
	defer controller.Finish()

	mockAuth := mocks.NewMockAuthService(controller)
	mockCore := mocks.NewMockCoreService(controller)
	mockService := service.Service{AuthService: mockAuth, CoreService: mockCore}

	header := NewHandler(config.Server{}, mockService)

	router := ginext.New("")
	router.GET("/events/:id", header.GetInfo)

	t.Run("invalid UUID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/events/"+invalidUUID, nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		require.Equal(t, http.StatusBadRequest, resp.Code)
		require.Contains(t, resp.Body.String(), errs.ErrInvalidEventID.Error())
	})

	t.Run("service GetInfo error", func(t *testing.T) {
		mockCore.EXPECT().GetInfo(gomock.Any(), testEventID).Return(nil, errors.New("internal error"))
		req := httptest.NewRequest(http.MethodGet, "/events/"+testEventID, nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		require.Equal(t, http.StatusInternalServerError, resp.Code)
	})

	t.Run("event not found", func(t *testing.T) {
		mockCore.EXPECT().GetInfo(gomock.Any(), testEventID).Return(nil, errs.ErrEventNotFound)
		req := httptest.NewRequest(http.MethodGet, "/events/"+testEventID, nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		require.Equal(t, http.StatusNotFound, resp.Code)
		require.Contains(t, resp.Body.String(), errs.ErrEventNotFound.Error())
	})

	t.Run("success", func(t *testing.T) {
		mockCore.EXPECT().GetInfo(gomock.Any(), testEventID).Return(&models.Event{
			Title:       "Test Event",
			Description: "Desc",
			Date:        time.Now().UTC(),
			Seats:       100,
		}, nil)
		req := httptest.NewRequest(http.MethodGet, "/events/"+testEventID, nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		require.Equal(t, http.StatusOK, resp.Code)
		require.Contains(t, resp.Body.String(), "Test Event")
	})

}

func TestHandler_CreateEvent(t *testing.T) {

	controller := gomock.NewController(t)
	defer controller.Finish()

	mockAuth := mocks.NewMockAuthService(controller)
	mockCore := mocks.NewMockCoreService(controller)
	mockService := service.Service{AuthService: mockAuth, CoreService: mockCore}

	header := NewHandler(config.Server{}, mockService)

	router := ginext.New("")
	router.POST("/events", header.CreateEvent)

	t.Run("no userID in context", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/events", bytes.NewBufferString(`{"title":"Test"}`))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		require.Equal(t, http.StatusUnauthorized, resp.Code)
		require.Contains(t, resp.Body.String(), errs.ErrInvalidToken.Error())
	})

	t.Run("invalid JSON binding", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), models.UserIDKey, testUserID)
		req := httptest.NewRequest(http.MethodPost, "/events", bytes.NewBufferString(`{"title":123}`))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(ctx)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		require.Equal(t, http.StatusBadRequest, resp.Code)
		require.Contains(t, resp.Body.String(), errs.ErrInvalidJSON.Error())
	})

	t.Run("missing date", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), models.UserIDKey, testUserID)
		req := httptest.NewRequest(http.MethodPost, "/events", bytes.NewBufferString(`{"title":"Test","description":"Desc","date":"","seats":10,"booking_ttl":"1h"}`))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(ctx)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		require.Equal(t, http.StatusBadRequest, resp.Code)
		require.Contains(t, resp.Body.String(), errs.ErrMissingDate.Error())
	})

	t.Run("invalid date", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), models.UserIDKey, testUserID)
		req := httptest.NewRequest(http.MethodPost, "/events", bytes.NewBufferString(`{"title":"Test","description":"Desc","date":"invalid","seats":10,"booking_ttl":"1h"}`))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(ctx)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		require.Equal(t, http.StatusBadRequest, resp.Code)
		require.Contains(t, resp.Body.String(), errs.ErrInvalidDate.Error())
	})

	t.Run("invalid booking TTL", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), models.UserIDKey, testUserID)
		req := httptest.NewRequest(http.MethodPost, "/events", bytes.NewBufferString(`{"title":"Test","description":"Desc","date":"2026-04-15T10:00:00Z","seats":10,"booking_ttl":"invalid"}`))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(ctx)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		require.Equal(t, http.StatusBadRequest, resp.Code)
		require.Contains(t, resp.Body.String(), errs.ErrInvalidBookingTTL.Error())
	})

	t.Run("service error", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), models.UserIDKey, testUserID)
		mockCore.EXPECT().CreateEvent(gomock.Any(), gomock.Any()).Return("", errors.New("internal error"))
		req := httptest.NewRequest(http.MethodPost, "/events", bytes.NewBufferString(`{"title":"Test","description":"Desc","date":"2026-04-15T10:00:00Z","seats":10,"booking_ttl":"1h"}`))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(ctx)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		require.Equal(t, http.StatusInternalServerError, resp.Code)
	})

	t.Run("success", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), models.UserIDKey, testUserID)
		mockCore.EXPECT().CreateEvent(gomock.Any(), gomock.Any()).Return("event123", nil)
		req := httptest.NewRequest(http.MethodPost, "/events", bytes.NewBufferString(`{"title":"Test","description":"Desc","date":"2026-04-15T10:00:00Z","seats":10,"booking_ttl":"1h"}`))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(ctx)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		require.Equal(t, http.StatusOK, resp.Code)
		require.Contains(t, resp.Body.String(), "event123")
	})

}

func TestHandler_CreateBooking(t *testing.T) {

	controller := gomock.NewController(t)
	defer controller.Finish()

	mockAuth := mocks.NewMockAuthService(controller)
	mockCore := mocks.NewMockCoreService(controller)
	mockService := service.Service{AuthService: mockAuth, CoreService: mockCore}

	header := NewHandler(config.Server{}, mockService)

	router := ginext.New("")
	router.POST("/events/:id/book", header.CreateBooking)

	t.Run("invalid UUID", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), models.UserIDKey, testUserID)
		req := httptest.NewRequest(http.MethodPost, "/events/"+invalidUUID+"/book", nil)
		req = req.WithContext(ctx)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		require.Equal(t, http.StatusBadRequest, resp.Code)
		require.Contains(t, resp.Body.String(), errs.ErrInvalidEventID.Error())
	})

	t.Run("no userID in context", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/events/"+testEventID+"/book", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		require.Equal(t, http.StatusUnauthorized, resp.Code)
		require.Contains(t, resp.Body.String(), errs.ErrInvalidToken.Error())
	})

	t.Run("service error", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), models.UserIDKey, testUserID)
		mockCore.EXPECT().CreateBooking(gomock.Any(), testUserID, testEventID).Return(int64(0), errors.New("internal error"))
		req := httptest.NewRequest(http.MethodPost, "/events/"+testEventID+"/book", nil)
		req = req.WithContext(ctx)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		require.Equal(t, http.StatusInternalServerError, resp.Code)
	})

	t.Run("success", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), models.UserIDKey, testUserID)
		mockCore.EXPECT().CreateBooking(gomock.Any(), testUserID, testEventID).Return(int64(456), nil)
		req := httptest.NewRequest(http.MethodPost, "/events/"+testEventID+"/book", nil)
		req = req.WithContext(ctx)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		require.Equal(t, http.StatusOK, resp.Code)
		require.Contains(t, resp.Body.String(), "456")
	})

}

func TestHandler_ConfirmBooking(t *testing.T) {

	controller := gomock.NewController(t)
	defer controller.Finish()

	mockAuth := mocks.NewMockAuthService(controller)
	mockCore := mocks.NewMockCoreService(controller)
	mockService := service.Service{AuthService: mockAuth, CoreService: mockCore}

	header := NewHandler(config.Server{}, mockService)

	router := ginext.New("")
	router.POST("/events/:id/confirm", header.ConfirmBooking)

	t.Run("invalid UUID", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), models.UserIDKey, testUserID)
		req := httptest.NewRequest(http.MethodPost, "/events/"+invalidUUID+"/confirm", nil)
		req = req.WithContext(ctx)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		require.Equal(t, http.StatusBadRequest, resp.Code)
		require.Contains(t, resp.Body.String(), errs.ErrInvalidEventID.Error())
	})

	t.Run("no userID in context", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/events/"+testEventID+"/confirm", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		require.Equal(t, http.StatusUnauthorized, resp.Code)
		require.Contains(t, resp.Body.String(), errs.ErrInvalidToken.Error())
	})

	t.Run("service error", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), models.UserIDKey, testUserID)
		mockCore.EXPECT().ConfirmBooking(gomock.Any(), testUserID, testEventID).Return(errors.New("internal error"))
		req := httptest.NewRequest(http.MethodPost, "/events/"+testEventID+"/confirm", nil)
		req = req.WithContext(ctx)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		require.Equal(t, http.StatusInternalServerError, resp.Code)
	})

	t.Run("success", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), models.UserIDKey, testUserID)
		mockCore.EXPECT().ConfirmBooking(gomock.Any(), testUserID, testEventID).Return(nil)
		req := httptest.NewRequest(http.MethodPost, "/events/"+testEventID+"/confirm", nil)
		req = req.WithContext(ctx)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		require.Equal(t, http.StatusOK, resp.Code)
		require.Contains(t, resp.Body.String(), string(models.StatusConfirmed))
	})

}
