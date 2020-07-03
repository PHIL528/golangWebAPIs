build:
	go build -o bin/server_exec server/main.go
	go build -o bin/listener_exec listener/main.go
	go build -o bin/client_exec client/main.go

clean:
	rm bin/*

auto: 
	go run script.go
m:
	bin/./server_exec
l: 
	bin/./listener_exec
c:
	bin/./client_exec make

pubsub:
	docker-compose up