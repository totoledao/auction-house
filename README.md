# auction-house

Real time auctioning!

docker compose up -d
go get -u ./...
go install github.com/jackc/tern/v2@latest
go run ./cmd/terndotenv
go run ./cmd/server
