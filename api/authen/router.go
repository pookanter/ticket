package authen

import (
	"ticket/pkg/apikit"
)

func Router(api *apikit.API) {
	h := NewHandler(api)

	api.App.POST("/sign-in", h.SignIn)
	api.App.POST("/sign-up", h.SignUp)

	api.App.POST("/refresh-token", h.RefreshToken)
}
