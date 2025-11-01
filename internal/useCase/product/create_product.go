package product

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/totoledao/auction-house/internal/validator"
)

type CreateProductReq struct {
	SellerID    uuid.UUID       `json:"seller_id"`
	ProductName string          `json:"product_name"`
	Description string          `json:"description"`
	BasePrice   decimal.Decimal `json:"base_price"`
	AuctionEnd  time.Time       `json:"auction_end"`
}

func (req CreateProductReq) Valid(ctx context.Context) validator.Evaluator {
	var eval validator.Evaluator
	const minAuctionDuration = 2 * time.Hour

	eval.CheckField(validator.NotBlank(req.ProductName), "product_name", "ProductName cannot be empty")
	eval.CheckField(validator.MaxChars(req.ProductName, 20), "product_name", "Product name too long")

	eval.CheckField(validator.NotBlank(req.Description), "description", "Description cannot be empty")
	eval.CheckField(validator.MinChars(req.Description, 10), "description", "Description too short")
	eval.CheckField(validator.MaxChars(req.Description, 200), "description", "Description too long")

	eval.CheckField(req.BasePrice.GreaterThan(decimal.NewFromInt(0)), "base_price", "BasePrice must be greater than 0")

	eval.CheckField(time.Until(req.AuctionEnd) >= minAuctionDuration, "auction_end", "AuctionEnd must be in 2 hours or more")

	return eval
}
