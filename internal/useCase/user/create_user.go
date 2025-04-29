package user

import (
	"context"

	"github.com/totoledao/auction-house/internal/validator"
)

type CreateUserReq struct {
	UserName string `json:"user_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Bio      string `json:"bio"`
}

func (req CreateUserReq) Valid(ctx context.Context) validator.Evaluator {
	var eval validator.Evaluator

	eval.CheckField(validator.NotBlank(req.UserName), "user_name", "Username cannot be empty")
	eval.CheckField(validator.MaxChars(req.UserName, 20), "user_name", "Username too long")

	eval.CheckField(validator.NotBlank(req.Email), "email", "Email cannot be empty")
	eval.CheckField(validator.Matches(req.Email, validator.EmailReg), "email", "Email is invalid")

	eval.CheckField(validator.MinChars(req.Password, 8), "password", "Password must have at least 8 characters")

	eval.CheckField(validator.NotBlank(req.Bio), "bio", "Bio cannot be empty")
	eval.CheckField(validator.MaxChars(req.Email, 500), "bio", "Bio too long")
	return eval
}
