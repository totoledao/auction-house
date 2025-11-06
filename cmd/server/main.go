package main

import (
	"context"
	"encoding/gob"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/totoledao/auction-house/internal/api"
	"github.com/totoledao/auction-house/internal/services"
)

func init() {
	gob.Register(uuid.UUID{})
}

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("POSTGRES_DB"),
	))
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	err = pool.Ping(ctx)
	if err != nil {
		panic(err)
	}

	s := scs.New()
	s.Store = pgxstore.New(pool)
	s.Lifetime = 24 * time.Hour
	s.Cookie.HttpOnly = true
	s.Cookie.SameSite = http.SameSiteLaxMode

	api := api.Api{
		Router:         chi.NewMux(),
		Sessions:       s,
		UserService:    services.NewUserService(pool),
		ProductService: services.NewProductService(pool),
		WsUpgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// DEV
				return true
			},
		},
		AuctionLobby: services.AuctionLobby{
			Rooms: make(map[uuid.UUID]*services.AuctionRoom),
		},
	}

	api.Routes()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hi"))
	})

	http.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Subroute"))
	})

	host := ":8080"
	fmt.Printf("Server running at http://localhost%s/ 🌐", host)
	err = http.ListenAndServe(host, api.Router)
	if err != nil {
		panic(err)
	}
}
