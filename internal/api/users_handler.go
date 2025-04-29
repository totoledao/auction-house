package api

import (
	"errors"
	"net/http"

	"github.com/totoledao/auction-house/internal/jsonutils"
	"github.com/totoledao/auction-house/internal/services"
	"github.com/totoledao/auction-house/internal/useCase/user"
)

func (api *Api) HandleSignUpUser(w http.ResponseWriter, r *http.Request) {
	data, problems, err := jsonutils.DecodeValidJson[user.CreateUserReq](r)
	if err != nil {
		jsonutils.EncodeJson(w, r, http.StatusUnprocessableEntity, problems)
		return
	}

	id, err := api.UserService.CreateUser(r.Context(), data.UserName, data.Email, data.Password, data.Bio)
	if err != nil {
		if errors.Is(err, services.ErrDuplicateEmailOrPassword) {
			_ = jsonutils.EncodeJson(w, r, http.StatusUnprocessableEntity, map[string]any{
				"error": services.ErrDuplicateEmailOrPassword.Error(),
			})
		}
		return
	}

	_ = jsonutils.EncodeJson(w, r, http.StatusCreated, map[string]any{
		"user_id": id,
	})
	return
}

func (api *Api) handleLoginUser(w http.ResponseWriter, r *http.Request)  {}
func (api *Api) handleLogoutUser(w http.ResponseWriter, r *http.Request) {}
