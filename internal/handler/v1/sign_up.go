package v1

import (
	"Kairos/internal/errs"
	"Kairos/internal/models"

	"github.com/wb-go/wbf/ginext"
)

// SignUp handles POST /api/v1/auth/sign-up.
// It binds the RegisterDTO, creates a new user via the service, generates a JWT token,
// and returns the token on success. Returns 409 if the login already exists.
func (h *Handler) SignUp(c *ginext.Context) {

	var request RegisterDTO
	if err := c.ShouldBindJSON(&request); err != nil {
		RespondError(c, errs.ErrInvalidJSON)
		return
	}

	userID, err := h.service.CreateUser(c.Request.Context(),
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
