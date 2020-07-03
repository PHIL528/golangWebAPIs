build:
	go build -o bin/client_exec client/main.go
	go build -o bin/server_exec server/main.go
	go build -o bin/listener_exec listener/main.go

pubsub:
	docker-compose up
