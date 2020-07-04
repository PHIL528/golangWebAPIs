FROM golang:1.12.0-alpine3.9
RUN apk add --no-cache git
WORKDIR /app/
COPY go.mod go.sum ./
RUN go mod download 
COPY . .
RUN go build -o bin/server_exec wait.go
CMD ["bin/./server_exec"]