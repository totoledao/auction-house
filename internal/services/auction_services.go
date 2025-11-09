package services

import (
	"context"
	"errors"
	"log/slog"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/shopspring/decimal"
)

type MessageKind int

const (
	// Request
	PlaceBid MessageKind = iota

	// Ok/Success
	SucceededToPlaceBid

	// Info
	NewBidPlaced
	AuctionFinished

	// Error
	FailedToPlaceBid
)

type Message struct {
	Message string
	Amount  decimal.Decimal
	Kind    MessageKind
	UserId  uuid.UUID
}

type AuctionLobby struct {
	sync.Mutex
	Rooms map[uuid.UUID]*AuctionRoom
}

type AuctionRoom struct {
	Id         uuid.UUID
	Context    context.Context
	Broadcast  chan Message
	Register   chan *Client
	Unregister chan *Client
	Clients    map[uuid.UUID]*Client

	BidService BidService
}

func (r *AuctionRoom) broadcastMessage(m Message) {
	slog.Info("New message received", "Message", m.Message, "RoomID", r.Id.String(), "UserID", m.UserId.String())
	switch m.Kind {
	case PlaceBid:
		bid, err := r.BidService.PlaceBid(r.Context, r.Id, m.UserId, m.Amount)

		client, ok := r.Clients[m.UserId]

		if err != nil {
			if errors.Is(err, BidErrors.MinValue) || errors.Is(err, BidErrors.PrevValue) {
				if ok {
					client.Send <- Message{Message: err.Error(), Kind: FailedToPlaceBid}
				}
				return
			}
		}

		if ok {
			client.Send <- Message{Message: "Your bid was successfully placed", Kind: SucceededToPlaceBid}
		}

		for id, client := range r.Clients {
			newMessage := Message{Message: "A new bid was placed", Kind: NewBidPlaced, Amount: bid.BidAmount}

			if id != m.UserId {
				client.Send <- newMessage
			}
		}
	}

}

func (r *AuctionRoom) registerClient(c *Client) {
	slog.Info("New user connected", "Client", c)
	r.Clients[*c.UserId] = c
}

func (r *AuctionRoom) unregisterClient(c *Client) {
	slog.Info("User disconnected", "Client", c)
	delete(r.Clients, *c.UserId)
}

func (r *AuctionRoom) Run() {
	defer func() {
		close(r.Broadcast)
		close(r.Register)
		close(r.Unregister)
	}()

	for {
		select {
		case message := <-r.Broadcast:
			r.broadcastMessage(message)
		case client := <-r.Register:
			r.registerClient(client)
		case client := <-r.Unregister:
			r.unregisterClient(client)
		case <-r.Context.Done():
			slog.Info("Auction finished", "AuctionID", r.Id)
			for _, client := range r.Clients {
				client.Send <- Message{Message: "Auction has been finished", Kind: AuctionFinished}
			}
			return
		}
	}
}

func NewAuctionRoom(ctx context.Context, id uuid.UUID, BidService BidService) *AuctionRoom {
	return &AuctionRoom{
		Id:         id,
		Context:    ctx,
		Broadcast:  make(chan Message),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		// Clients:    map[uuid.UUID]*Client,
		BidService: BidService,
	}
}

type Client struct {
	Room   *AuctionRoom
	Conn   *websocket.Conn
	Send   chan Message
	UserId *uuid.UUID
}

func NewClient(
	room *AuctionRoom,
	conn *websocket.Conn,
	userId *uuid.UUID,
) *Client {
	return &Client{
		Room:   room,
		Conn:   conn,
		Send:   make(chan Message, 512),
		UserId: userId,
	}
}
