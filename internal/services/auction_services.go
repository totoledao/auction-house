package services

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"sync"
	"time"

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
	InvalidJSON
)

type Message struct {
	Message string          `json:"message,omitempty"`
	Amount  decimal.Decimal `json:"amount"`
	Kind    MessageKind     `json:"kind"`
	UserId  uuid.UUID       `json:"user_id,omitempty"`
}

type AuctionLobby struct {
	sync.Mutex
	Rooms map[uuid.UUID]*AuctionRoom
}

func (l *AuctionLobby) AssignProductToRoom(id uuid.UUID, auctionEnd time.Time, BidService BidService) {
	ctx, _ := context.WithDeadline(context.Background(), auctionEnd)

	auctionRoom := newAuctionRoom(ctx, id, BidService)

	go auctionRoom.Run()

	l.Lock()
	l.Rooms[id] = auctionRoom
	l.Unlock()
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
		bid, err := r.BidService.PlaceBid(r.Context, m.UserId, r.Id, m.Amount)

		client, ok := r.Clients[m.UserId]

		if err != nil {
			slog.Error("Error placing bid", "err", err)

			if errors.Is(err, BidErrors.MinValue) || errors.Is(err, BidErrors.PrevValue) {
				if ok {
					client.Send <- Message{Message: err.Error(), Kind: FailedToPlaceBid, Amount: m.Amount, UserId: m.UserId}
				}
				return
			}

			if ok {
				client.Send <- Message{Message: "Something went wrong, try again", Kind: FailedToPlaceBid, Amount: m.Amount, UserId: m.UserId}
			}
			return
		}

		if ok {
			client.Send <- Message{Message: "Your bid was successfully placed", Kind: SucceededToPlaceBid, Amount: m.Amount, UserId: m.UserId}
		}

		for id, client := range r.Clients {
			newMessage := Message{Message: "A new bid was placed", Kind: NewBidPlaced, Amount: bid.BidAmount, UserId: m.UserId}

			if id != m.UserId {
				client.Send <- newMessage
			}
		}

	case InvalidJSON:
		client, ok := r.Clients[m.UserId]
		if !ok {
			slog.Info("Client not found in hashmap", "userId", m.UserId)
			return
		}
		client.Send <- m
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
	slog.Info("Auction has begun", "auctionId", r.Id)
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

func newAuctionRoom(ctx context.Context, id uuid.UUID, BidService BidService) *AuctionRoom {
	return &AuctionRoom{
		Id:         id,
		Context:    ctx,
		Broadcast:  make(chan Message),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[uuid.UUID]*Client),
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

const (
	maxMessageSize = 512                     // 512 bytes
	readDeadline   = 60 * time.Second        // 60 seconds
	pingPeriod     = (readDeadline * 9) / 10 // 90% of readDeadline
	writeWait      = 10 * time.Second        // 10 seconds
)

func (c *Client) ReadEventLoop() {
	defer func() {
		c.Room.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(readDeadline))
	c.Conn.SetPongHandler(func(appData string) error {
		c.Conn.SetReadDeadline(time.Now().Add(readDeadline))
		return nil
	})

	for {
		_, data, err := c.Conn.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				slog.Error("Unexpected ws close error", "error", err)
			}
			return
		}

		var m Message
		err = json.Unmarshal(data, &m)
		if err != nil {
			c.Room.Broadcast <- Message{
				Kind:    InvalidJSON,
				Message: "This message should be valid JSON",
				UserId:  *c.UserId,
			}
			continue
		}

		m.UserId = *c.UserId

		c.Room.Broadcast <- m
	}
}

func (c *Client) WriteEventLoop() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Conn.WriteJSON(Message{
					Message: "closing ws connection",
					Kind:    websocket.CloseMessage,
				})
				return
			}

			if message.Kind == AuctionFinished {
				close(c.Send)
				return
			}
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))

			err := c.Conn.WriteJSON(message)
			if err != nil {
				c.Room.Unregister <- c
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))

			err := c.Conn.WriteMessage(websocket.PingMessage, nil)
			if err != nil {
				slog.Error("Unexpected write error", "error", err)
				return
			}
		}
	}
}
