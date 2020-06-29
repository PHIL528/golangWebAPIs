package main

import (
	"log"
	"os"
	"proto-playground/proto"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {
	client_name := os.Args[1]

	var con *grpc.ClientConn
	con, err := grpc.Dial("localhost:3001", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Error %v", err)
	}
	defer con.Close()
	client := proto.NewReservationServiceClient(con)

	book_trip_request := proto.BookTrip{
		PassengerName: client_name,
	}

	confirmed_trip, err := client.MakeReservation(context.Background(), &book_trip_request)
	if err != nil {
		log.Fatalf("Error %v", err)
	}
	log.Printf("Server assigned driver %s", confirmed_trip.DriverName)

}
