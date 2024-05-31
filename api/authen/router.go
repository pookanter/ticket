package authen

import (
	"ticket/api/authen/authorize"
	"ticket/api/authen/users"
	"ticket/pkg/apikit"
	"ticket/pkg/auth"
)

func Router(api *apikit.API) {
	a := authorize.New(api)

	api.App.POST("/sign-in", a.SignIn)
	api.App.POST("/sign-up", a.SignUp)
	api.App.POST("/refresh-token", a.RefreshToken)

	u := users.New(api)

	usersGroup := api.App.Group("/users")
	usersGroup.Use(auth.Middleware(api.Config))
	usersGroup.GET("/me", u.GetMe)
}
