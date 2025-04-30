package user

import (
	"context"

	"github.com/totoledao/auction-house/internal/validator"
)

type LoginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (req LoginUserRequest) Valid(ctx context.Context) validator.Evaluator {
	var eval validator.Evaluator

	eval.CheckField(validator.Matches(req.Email, validator.EmailReg), "email", "invalid email")
	eval.CheckField(validator.NotBlank(req.Password), "password", "password cannot be empty")

	return eval
}
