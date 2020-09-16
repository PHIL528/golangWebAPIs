# proto-playground

## Simple REST, gRPC and Pub/Sub testing

<br></br>

During the first two weeks of my internship at Karhoo, my mentor Marek Chmiel had me work on this repository. There is a client and server terminal, where the client sends a request to the server and the server sends back a response. 

The client can send data via REST, gRPC, or Pub/Sub, and the server listens with a thread for each of the three routes. 

Uses Maxbrunsfield's Counterfieter library is used to generate mocks for each API type. Stub responses are used for gRPC, Golang's httptest package is used for REST, and recorded innvocations are used to mock the Pub/Sub as a function of the input. 

Requires Google Cloud Pub/Sub to be installed locally to handle Pub/Sub. 




**Installation** 
 
```cd go/src ```

```git clone https://github.com/PHIL528/golangWebAPIs```

<div> <h1></h1></div>

 
**Running the Emulator, Server, Client and Listener**

Server:

```cd go/src/proto-playground/server```

```go run main.go```

Client:

```cd go/src/proto-playground/client```

```go run main.go <rest OR grpc OR pubsub> <firstname>```

Then click on the stock terminal icon, and there should be 4 terminals

Find the terminal for the client, and enter in either <rest> <grpc> or <pubsub> and then <firstname>


