package rest

import (
	"auth/internal/constants"
	"auth/internal/domain/models"
	"auth/internal/ports"
	"auth/internal/utils"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type Handler struct {
	AuthService ports.AuthService
	Log         *logrus.Logger
	helpers     *helpers
}

func NewHandler(authService ports.AuthService, log *logrus.Logger, helpers *helpers) *Handler {
	return &Handler{AuthService: authService,
		Log:     log,
		helpers: helpers,
	}
}

//func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
//	ctx, span := utils.StartSpan(r.Context())
//	defer span.End()
//
//	var userData models.User
//
//	login, password, ok := r.BasicAuth()
//	if !ok {
//		h.helpers.Error(w, r, ErrorBadRequest(fmt.Errorf("error parsing basic auth data")))
//		return
//	}
//
//	userData.Login = login
//	userData.Password = password
//
//	if err := h.AuthService.SignUp(ctx, &userData); err != nil {
//		if errors.Is(err, models.ErrUserExist) {
//			h.helpers.Error(w, r, ErrConflict(err))
//			return
//		}
//		h.helpers.Error(w, r, ErrorInternal(err))
//		return
//	}
//
//	h.Log.Println("registered new user")
//	http.Redirect(w, r, "/auth/login", http.StatusOK)
//
//}

// Login
// @ID login
// @tags auth
// @Summary Authorized user
// @Description Authenticate and authorized user. Return access and refresh tokens in cookies.
// @Accept json
// @Produce json
// @Param redirect_uri query string false "redirect uri"
// @Param Login body models.User true "request body"
// @Success 200 {object} models.TokenPair "ok"
// @Header 200 {string} accessToken	"token for access services"
// @Header 200 {string} refreshToken "token for refresh access_token"
// @Failure 400 {object} rest.Error "bad request"
// @Failure 404 {string} string "404 page not found"
// @Failure 403 {object} rest.Error "forbidden"
// @Failure 500 {object} rest.Error "internal error"
// @Router /login [post]
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	ctx, span := utils.StartSpan(r.Context())
	defer span.End()

	var reqData models.User
	//var tokens *models.TokenDetails
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&reqData)
	if err != nil {
		h.helpers.Error(w, r, ErrorInternal(err))
		return
	}

	tokens, err := h.AuthService.Authorize(ctx, &reqData)
	if err != nil {
		h.helpers.Error(w, r, ErrorForbidden(err))
		return
	}

	atCookie := http.Cookie{
		Name:     constants.ACCESS_TOKEN,
		Value:    tokens.AccessToken,
		Path:     "/",
		Expires:  tokens.AtExpires,
		HttpOnly: true,
	}

	rtCookie := http.Cookie{
		Name:     constants.REFRESH_TOKEN,
		Value:    tokens.RefreshToken,
		Path:     "/",
		Expires:  tokens.RtExpires,
		HttpOnly: true,
	}

	http.SetCookie(w, &atCookie)
	http.SetCookie(w, &rtCookie)

	redirectUrl := r.URL.Query().Get(constants.REDIRECT_URI)

	if len(redirectUrl) > 0 {
		http.Redirect(w, r, redirectUrl, http.StatusFound)
	} else {
		h.helpers.ResponseOk(w, r, tokens)
	}
}

// Logout
// @ID logout
// @tags auth
// @Summary Clears tokens
// @Description Clears access and refresh tokens
// @Security Auth
// @Produce json
// @Param redirect_uri query string false "redirect uri"
// @Param accessToken header string true "access token"
// @Param refreshToken header string true "refresh token"
// @Success 200  "ok"
// @Failure 302  "redirect"
// @Failure 500  "internal error"
// @Router /logout [post]
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	access := http.Cookie{
		Name:     "access",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	}
	refresh := http.Cookie{
		Name:     "refresh",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	}
	http.SetCookie(w, &access)
	http.SetCookie(w, &refresh)

	redirectUrl := r.URL.Query().Get(constants.REDIRECT_URI)

	if len(redirectUrl) > 0 {
		http.Redirect(w, r, redirectUrl, http.StatusFound)
	} else {
		w.WriteHeader(http.StatusOK)
	}

}

// Info
// @ID Info
// @tags auth
// @Summary Validate tokens
// @Description Validate tokens and refresh tokens if refresh token is valid
// @Security Auth
// @Produce json
// @Param accessToken header string true "access token"
// @Param refreshToken header string true "refresh token"
// @Success 200 {object} models.User "ok"
// @Failure 403 {string} string "forbidden"
// @Failure 500 {string} string "internal error"
// @Router /i [post]
func (h *Handler) Info(w http.ResponseWriter, r *http.Request) {
	ctx, span := utils.StartSpan(r.Context())
	defer span.End()

	access, err := r.Cookie("access")
	if err != nil {
		h.helpers.Error(w, r, ErrorForbidden(err))
		return
	}

	refresh, err := r.Cookie("refresh")
	if err != nil {
		h.helpers.Error(w, r, ErrorForbidden(err))
		return
	}

	tokens := &models.TokenPair{
		AccessToken:  access.Value,
		RefreshToken: refresh.Value,
	}

	user, err := h.AuthService.ValidateTokens(ctx, tokens)
	if err != nil {
		h.helpers.Error(w, r, ErrorForbidden(err))
		return
	}
	h.helpers.ResponseOk(w, r, user.Login)
}
