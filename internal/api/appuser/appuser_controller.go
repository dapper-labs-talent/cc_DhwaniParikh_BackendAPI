package appuser

import (
	"fmt"
	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"
	"net/http"
	"server/internal/auth"
	"server/internal/models"
	"server/internal/repositories"
	"server/internal/util/response"
)

type AppUserController struct {
	repo repositories.IAppUserRepository
}

func NewAppUserController(r *repositories.AppUserRepo) *AppUserController {
	return &AppUserController{
		repo: r,
	}
}

func (ac AppUserController) Routes() []models.Route {
	return []models.Route{
		{
			Method:  http.MethodGet,
			Path:    "/",
			Handler: ac.GetAllUsers,
		},
		{
			Method:  http.MethodPut,
			Path:    "/",
			Handler: ac.UpdateUser,
		},
	}
}

// GetAllUsers Returns all user details
// @Summary Get all users
// @Description Returns information for all user account
// @Accept json
// @Produce json
// @Success 200 {object} []models.AppUserEntities User Account Information
// @Failure 400 {object} response.ErrResponse Not found
// @Failure 401 {object} response.ErrResponse Not found
// @Failure 404 {object} response.ErrResponse Not found
// @Failure 500 {object} response.ErrResponse Internal Error
// @Router /api/users [get]
func (ac *AppUserController) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	fmt.Println("test")
	users, err := ac.repo.GetAllUsers(ctx)
	if len(users) == 0 {
		response.NotFound(w, r)
		return
	}
	if err != nil {
		log.Ctx(ctx).Err(err).Msg(fmt.Sprintf("failed to get user "))
		response.InternalServerError(w, r, response.ErrResponse{
			Err:        err,
			StatusText: appCodeText(appCodeAccountLoadError),
			AppCode:    appCodeAccountLoadError,
		})
		return
	}

	response.Ok(w, r, users)
}

// UpdateUser updates user info
// @Updates a specific user account
// @Accept json
// @Produce json
// @Success 200 {object} Update user details
// @Failure 400 {object} response.ErrResponse Not found
// @Failure 401 {object} response.ErrResponse Not found
// @Failure 404 {object} response.ErrResponse Not found
// @Failure 500 {object} response.ErrResponse Internal Error
// @Router /api/users [put]
func (ac *AppUserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	_, claims, _ := auth.FromContext(r.Context())
	userID := claims.Id

	ur := new(models.UserUpdateRequest)
	if err := render.Bind(r, ur); err != nil {
		log.Ctx(ctx).Err(err).Msg("failed to bind Data")
		response.BadRequest(w, r, response.ErrResponse{
			Err:        err,
			StatusText: appCodeText(appCodeBindFailed),
			AppCode:    appCodeBindFailed,
		})

		return
	}

	account, err := ac.repo.GetById(ctx, userID)
	if account == nil {
		response.NotFound(w, r)
		return
	}
	if err != nil {
		log.Ctx(ctx).Err(err).Msg(fmt.Sprintf("failed to get user with id: %s", userID.String()))
		response.InternalServerError(w, r, response.ErrResponse{
			Err:        err,
			StatusText: appCodeText(appCodeAccountLoadError),
			AppCode:    appCodeAccountLoadError,
		})
		return
	}

	account.FirstName = ur.FirstName
	account.LastName = ur.LastName
	err = nil
	err = ac.repo.Update(ctx, account)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg(fmt.Sprintf("unable to update account for user with id: %s", account.ID.String()))
		response.InternalServerError(w, r, response.ErrResponse{
			Err:        err,
			StatusText: appCodeText(appCodeAccountUpdateError),
			AppCode:    appCodeAccountUpdateError,
		})
		return
	}

	response.Ok(w, r, nil)
}
