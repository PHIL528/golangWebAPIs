build:
    go build -o bin/main client/main.go
    go build -o bin/main server/main.go
    go build -o bin/main listener/main.go
    go build -o bin/main Config/Config.go


run:
    go run server/main.go
    go run bin/main
    go run client.go
    go run listener.go

pubsub: 
	docker-compose up
