// Package v1 provides version 1 of the Kairos API handlers for notifications.
// It includes endpoints to create, query, and cancel notifications via HTTP.
package v1

import (
	"Kairos/internal/errs"
	"Kairos/internal/service"

	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/helpers"
)

// Handler is the v1 API handler for notifications.
// It wraps the service layer and provides HTTP endpoints for CRUD operations.
type Handler struct {
	service service.Service
}

// NewHandler creates a new v1 Handler with the provided service.
func NewHandler(service service.Service) *Handler {
	return &Handler{service: service}
}

// GetNotification handles GET /notify?id=<id> requests.
// It validates the notification ID and returns the current status of the notification.
// Returns an error if the ID is invalid or the notification is not found.
func (h *Handler) GetNotification(c *ginext.Context) {

	notificationID := c.Query("id")
	if err := helpers.ParseUUID(notificationID); err != nil {
		RespondError(c, errs.ErrInvalidNotificationID)
		return
	}

	status, err := h.service.GetStatus(c.Request.Context(), notificationID)
	if err != nil {
		RespondError(c, err)
		return
	}

	respondOK(c, status)

}

// CancelNotification handles DELETE /notify?id=<id> requests.
// It validates the notification ID and cancels the notification if possible.
// Returns an error if the ID is invalid or the notification cannot be canceled.
func (h *Handler) CancelNotification(c *ginext.Context) {

	notificationID := c.Query("id")
	if err := helpers.ParseUUID(notificationID); err != nil {
		RespondError(c, errs.ErrInvalidNotificationID)
		return
	}

	if err := h.service.CancelNotification(c.Request.Context(), notificationID); err != nil {
		RespondError(c, err)
		return
	}

	respondOK(c, "canceled")

}
