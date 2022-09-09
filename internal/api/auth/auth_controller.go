package auth

import (
	"fmt"
	"github.com/go-chi/render"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"server/internal/auth"
	"server/internal/config"
	"server/internal/entities"
	"server/internal/models"
	"server/internal/repositories"
	"server/internal/util/response"
)

type AuthController struct {
	appUserRepo repositories.IAppUserRepository
	token       *auth.TokenService
	config      *config.Config
	logger      *zerolog.Logger
}

func NewAuthController(c *config.Config, r repositories.IAppUserRepository, t *auth.TokenService, l *zerolog.Logger) *AuthController {
	return &AuthController{
		appUserRepo: r,
		token:       t,
		config:      c,
		logger:      l,
	}
}

func (ac AuthController) Routes() []models.Route {
	return []models.Route{
		{
			Method:  http.MethodPost,
			Path:    "/login",
			Handler: ac.Login,
		},
		{
			Method:  http.MethodPost,
			Path:    "/signup",
			Handler: ac.Signup,
		},
	}
}

// Login user
// @Summary Login user
// @Accept json
// @Produce json
// @Success 201 {object} Token
// @Failure 400 {object} response.ErrResponse Not found
// @Failure 401 {object} response.ErrResponse Not found
// @Failure 404 {object} response.ErrResponse Not found
// @Failure 500 {object} response.ErrResponse Internal Error
// @Router /auth/signup [post]
func (ac AuthController) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	lr := new(models.UserLoginRequest)

	if err := render.Bind(r, lr); err != nil {
		log.Ctx(ctx).Err(err).Msg("Unable to bind request AdminLoginRequest")
		response.BadRequest(w, r, response.ErrResponse{
			Err:        err,
			StatusText: appCodeText(appCodeBindFailed),
			AppCode:    appCodeBindFailed,
		})

		return
	}

	user, err := ac.appUserRepo.GetByEmail(ctx, lr.Email)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg(fmt.Sprintf("Unable to get admin user by email: %s", lr.Email))
		response.InternalServerError(w, r, response.ErrResponse{
			Err:        err,
			StatusText: appCodeText(appCodeAccountLoadError),
			AppCode:    appCodeAccountLoadError,
		})
		return
	}

	if user == nil {
		log.Ctx(ctx).Err(err).Msg(fmt.Sprintf("User not found"))
		response.InternalServerError(w, r, response.ErrResponse{
			Err:        err,
			StatusText: appCodeText(appCodeAccountLoadError),
			AppCode:    appCodeAccountLoadError,
		})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(lr.Password))
	if err != nil {
		log.Ctx(ctx).Err(err).Msg(fmt.Sprintf("Wrong Credentials"))
		response.InternalServerError(w, r, response.ErrResponse{
			Err:        err,
			StatusText: appCodeText(appCodeAccountLoadError),
			AppCode:    appCodeAccountLoadError,
		})
		return
	}

	tokens, err := ac.token.GenerateAPITokens(user.ID, lr.Email)

	if err != nil {
		log.Ctx(ctx).Err(err).Msg(fmt.Sprintf("Unable to generate admin tokens for user with id: %s", user.ID.String()))
		response.InternalServerError(w, r, response.ErrResponse{
			Err:        err,
			StatusText: appCodeText(appCodeTokenError),
			AppCode:    appCodeTokenError,
		})
		return
	}
	response.Ok(w, r, tokens)
}

// Signup a new user account
// @Summary Creates account
// @Description Creates account
// @Accept json
// @Produce json
// @Param create body models.AppUserEntity true "Create Request"
// @Success 201 {object} Token
// @Failure 400 {object} response.ErrResponse Not found
// @Failure 401 {object} response.ErrResponse Not found
// @Failure 404 {object} response.ErrResponse Not found
// @Failure 500 {object} response.ErrResponse Internal Error
// @Router /auth/signup [post]
func (ac AuthController) Signup(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	lr := new(models.UserSignupRequest)

	if err := render.Bind(r, lr); err != nil {
		log.Ctx(ctx).Err(err).Msg("Unable to bind request AdminLoginRequest")
		response.BadRequest(w, r, response.ErrResponse{
			Err:        err,
			StatusText: appCodeText(appCodeBindFailed),
			AppCode:    appCodeBindFailed,
		})

		return
	}
	pb, gErr := bcrypt.GenerateFromPassword([]byte(lr.Password), bcrypt.DefaultCost)
	ph := string(pb)
	if gErr != nil {
		log.Ctx(ctx).Err(gErr).Msg(fmt.Sprintf("unable to create password for account "))
		response.InternalServerError(w, r, response.ErrResponse{
			Err:        gErr,
			StatusText: appCodeText(appCodeAccountUpdateError),
			AppCode:    appCodeAccountUpdateError,
		})
		return
	}
	account := &entities.AppUserEntity{
		ID:        uuid.NewV4(),
		Email:     lr.Email,
		FirstName: lr.FirstName,
		LastName:  lr.LastName,
		Password:  ph,
	}

	user, err := ac.appUserRepo.Create(ctx, account)
	if err != nil {
		fmt.Println(err.Error())
		log.Ctx(ctx).Err(err).Msg("unable for create account")
		response.InternalServerError(w, r, response.ErrResponse{
			Err:        err,
			StatusText: appCodeText(appCodeCreateError),
			AppCode:    appCodeCreateError,
		})
		return
	}

	if user == nil {
		response.Unauthorized(w, r)
		return
	}

	tokens, err := ac.token.GenerateAPITokens(user.ID, lr.Email)

	if err != nil {
		log.Ctx(ctx).Err(err).Msg(fmt.Sprintf("Unable to generate admin tokens for user with id: %s", user.ID.String()))
		response.InternalServerError(w, r, response.ErrResponse{
			Err:        err,
			StatusText: appCodeText(appCodeTokenError),
			AppCode:    appCodeTokenError,
		})
		return
	}
	response.Ok(w, r, tokens)
}
