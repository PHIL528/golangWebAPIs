# proto-playground

## Simple "hello world" gRPC test

**Installation** 
 
```cd go/src ```

```git clone https://github.com/marchmiel/proto-playground```

<div> <h1></h1></div>
 
**Running the Server**

```cd go/src/proto-playground/server```

```go run main.go```
<div> <h1></h1></div>

**Running the Client**

```cd go/src/proto-playground/client```

```go run main.go <first_name>```


I watched two videos on gRPC, both of them used .WithInsecure(), and one of them used the relfection package which I guess is not required for a hello-world example but is good practice.

As you per the diagram you wrote, the client sends a BookTrip message which stores the passenger name, and the server sends a Trip message which stores passanger name and driver name.


