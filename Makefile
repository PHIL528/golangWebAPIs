build:
    go build -o client/main.go
    go build -o server/main.go
    go build -o listener/main.go
    go build -o Config/Config.go


run:
	docker-compose up
    go run server/main.go
    go run client/main.go
    go run listener/main.go


