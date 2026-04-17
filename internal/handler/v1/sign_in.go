package v1

import (
	"Kairos/internal/errs"
	"Kairos/internal/models"

	"github.com/wb-go/wbf/ginext"
)

// SignIn handles POST /api/v1/auth/sign-in.
// It binds the LoginDTO, authenticates the user via the service, generates a JWT token,
// and returns the token on success. Returns 401 for invalid credentials.
func (h *Handler) SignIn(c *ginext.Context) {

	var request LoginDTO
	if err := c.ShouldBindJSON(&request); err != nil {
		RespondError(c, errs.ErrInvalidJSON)
		return
	}

	userID, err := h.service.GetUserId(c.Request.Context(),
		models.User{Login: request.Login, Password: request.Password})
	if err != nil {
		RespondError(c, err)
		return
	}

	token, err := h.service.CreateToken(userID)
	if err != nil {
		RespondError(c, err)
		return
	}

	respondOK(c, token)

}
