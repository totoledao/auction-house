package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
	"github.com/totoledao/auction-house/internal/store/pgstore"
)

type ProductService struct {
	pool    *pgxpool.Pool
	queries *pgstore.Queries
}

func NewProductService(pool *pgxpool.Pool) ProductService {
	return ProductService{
		pool:    pool,
		queries: pgstore.New(pool),
	}
}

func (ps *ProductService) CreateProduct(
	ctx context.Context,
	seller_id uuid.UUID,
	product_name string,
	description string,
	base_price decimal.Decimal,
	auction_end time.Time,
) (uuid.UUID, error) {
	id, err := ps.queries.CreateProduct(
		ctx,
		pgstore.CreateProductParams{
			SellerID:    seller_id,
			ProductName: product_name,
			Description: description,
			BasePrice:   base_price,
			AuctionEnd:  auction_end,
		},
	)

	if err != nil {
		return uuid.UUID{}, err
	}
	return id, nil
}
