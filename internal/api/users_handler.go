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
		if errors.Is(err, services.ErrDuplicateEmailOrUserName) {
			_ = jsonutils.EncodeJson(w, r, http.StatusUnprocessableEntity, map[string]any{
				"error": services.ErrDuplicateEmailOrUserName.Error(),
			})
		}
		return
	}

	_ = jsonutils.EncodeJson(w, r, http.StatusCreated, map[string]any{
		"user_id": id,
	})
}

func (api *Api) handleLoginUser(w http.ResponseWriter, r *http.Request) {
	data, problems, err := jsonutils.DecodeValidJson[user.LoginUserRequest](r)
	if err != nil {
		jsonutils.EncodeJson(w, r, http.StatusUnprocessableEntity, problems)
		return
	}

	id, err := api.UserService.AuthenticateUser(r.Context(), data.Email, data.Password)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			_ = jsonutils.EncodeJson(w, r, http.StatusBadRequest, map[string]any{
				"error": services.ErrInvalidCredentials.Error(),
			})
			return
		}

		_ = jsonutils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{
			"error": "unexpected internal server error",
		})
		return
	}

	err = api.Sessions.RenewToken(r.Context())
	if err != nil {
		_ = jsonutils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{
			"error": "unexpected internal server error",
		})
		return
	}

	api.Sessions.Put(r.Context(), "AuthenticatedUserId", id)
	_ = jsonutils.EncodeJson(w, r, http.StatusOK, map[string]any{
		"message": "logged in successfully",
	})
}

func (api *Api) handleLogoutUser(w http.ResponseWriter, r *http.Request) {
	err := api.Sessions.RenewToken(r.Context())
	if err != nil {
		_ = jsonutils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{
			"error": "unexpected internal server error",
		})
		return
	}

	api.Sessions.Remove(r.Context(), "AuthenticatedUserId")
	_ = jsonutils.EncodeJson(w, r, http.StatusOK, map[string]any{
		"message": "logged out successfully",
	})
}
