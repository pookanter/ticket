package users

import (
	"net/http"
	"ticket/pkg/apikit"
	"ticket/pkg/auth"
	"ticket/pkg/db"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	App      *echo.Echo
	Auth     *auth.Auth
	DB       *db.Queries
	DBConfig apikit.DBConfig
}

func NewHandler(api *apikit.API) *Handler {
	return &Handler{
		App: api.App,
		Auth: auth.New(auth.AuthConfig{
			RSAKey:             api.GetCerts().PrivateKey,
			AccessTokenExpire:  api.GetGlobalConfig().AccessTokenExpire,
			RefreshTokenExpire: api.GetGlobalConfig().RefreshTokenExpire,
		}),
		DB:       api.DB,
		DBConfig: api.GetDBConfig(),
	}
}

func (h *Handler) GetMe(c echo.Context) error {
	ac := c.(auth.AuthContext)
	user, err := h.DB.FindUserByID(c.Request().Context(), ac.Claims.UserID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if user.ID == 0 {
		return echo.NewHTTPError(http.StatusNotFound, "user not found")
	}

	return c.JSON(http.StatusOK, struct {
		ID        uint64 `json:"id"`
		Email     string `json:"email"`
		Name      string `json:"name"`
		Lastname  string `json:"lastname"`
		CreatedAt string `json:"created_at"`
	}{
		ID:        user.ID,
		Email:     user.Email.String,
		Name:      user.Name.String,
		Lastname:  user.Lastname.String,
		CreatedAt: user.CreatedAt.Time.Format("2006-01-02 15:04:05"),
	})
}
