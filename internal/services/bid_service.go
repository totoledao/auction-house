package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
	"github.com/totoledao/auction-house/internal/store/pgstore"
)

type BidService struct {
	pool    *pgxpool.Pool
	queries *pgstore.Queries
}

func NewBidService(pool *pgxpool.Pool) BidService {
	return BidService{
		pool:    pool,
		queries: pgstore.New(pool),
	}
}

type errBid struct {
	MinValue  error
	PrevValue error
}

var BidErrors = errBid{
	MinValue:  errors.New("the bid value should be higher than the minimum value"),
	PrevValue: errors.New("the bid should be higher than the previous bid"),
}

func (ps *BidService) PlaceBid(
	ctx context.Context,
	bidder_id uuid.UUID,
	product_id uuid.UUID,
	bid_amount decimal.Decimal,
) (pgstore.Bid, error) {
	product, err := ps.queries.GetProductById(
		ctx,
		product_id,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pgstore.Bid{}, err
		}
	}
	if product.BasePrice.GreaterThanOrEqual(bid_amount) {
		return pgstore.Bid{}, BidErrors.MinValue
	}

	highestBid, err := ps.queries.GetHighestBidByProductId(
		ctx,
		product_id,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pgstore.Bid{}, err
		}
	}
	if highestBid.BidAmount.GreaterThanOrEqual(bid_amount) {
		return pgstore.Bid{}, BidErrors.PrevValue
	}

	data, err := ps.queries.CreateBid(
		ctx,
		pgstore.CreateBidParams{
			ProductID: product_id,
			BidderID:  bidder_id,
			BidAmount: bid_amount,
		},
	)
	if err != nil {
		return pgstore.Bid{}, err
	}
	return data, nil
}
