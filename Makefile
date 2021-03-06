build:
	go build -o bin/server_exec server/main.go
	go build -o bin/listener_exec listener/main.go
	go build -o bin/client_exec client/main.go

clean:
	rm bin/*

run: 
	go run script.go
s:
	bin/./server_exec
l: 
	bin/./listener_exec
c:
	bin/./client_exec make

pubsub:
	docker-compose up

cgen:
	go run github.com/maxbrunsfeld/counterfeiter/v6 ./client/main.go PubSubConnector