package api

import (
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"github.com/totoledao/auction-house/internal/services"
)

type Api struct {
	Router         *chi.Mux
	Sessions       *scs.SessionManager
	UserService    services.UserService
	ProductService services.ProductService
	BidService     services.BidService
	WsUpgrader     websocket.Upgrader
	AuctionLobby   services.AuctionLobby
}
