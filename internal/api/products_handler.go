package api

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/totoledao/auction-house/internal/jsonutils"
	"github.com/totoledao/auction-house/internal/services"
	"github.com/totoledao/auction-house/internal/useCase/product"
)

func (api *Api) HandleCreateProduct(w http.ResponseWriter, r *http.Request) {
	data, problems, err := jsonutils.DecodeValidJson[product.CreateProductReq](r)

	if err != nil {
		jsonutils.EncodeJson(w, r, http.StatusUnprocessableEntity, problems)
		return
	}

	userID, ok := api.Sessions.Get(r.Context(), "AuthenticatedUserId").(uuid.UUID)
	if !ok {
		jsonutils.EncodeJson(w, r, http.StatusInternalServerError,
			map[string]any{
				"error": "Unexpected error. Try again later.",
			})
		return
	}

	productId, err := api.ProductService.CreateProduct(r.Context(),
		userID,
		data.ProductName,
		data.Description,
		data.BasePrice,
		data.AuctionEnd,
	)
	if err != nil {
		jsonutils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{
			"error": "Failed to create product auction. Try again later.",
		})
		return
	}

	ctx, _ := context.WithDeadline(context.Background(), data.AuctionEnd)

	auctionRoom := services.NewAuctionRoom(ctx, productId, api.BidService)

	go auctionRoom.Run()

	api.AuctionLobby.Lock()
	api.AuctionLobby.Rooms[productId] = auctionRoom
	api.AuctionLobby.Unlock()

	jsonutils.EncodeJson(w, r, http.StatusCreated, map[string]any{
		"message":    "Product auction started",
		"product_id": productId,
	})
}
