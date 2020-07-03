build:
	go build -o bin/clientexec client/main.go
	go build -o bin/serverexec server/main.go
	go build -o bin/listenerexec listener/main.go

pubsub:
	docker-compose up
