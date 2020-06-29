package main

import (
	"context"
	"log"
	"net"
	"proto-playground/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	lis, err := net.Listen("tcp", ":3001") //React port + 1
	if err != nil {
		log.Fatalf("Error %v", err)
	}
	gcrpServer := grpc.NewServer()
	//proto generates ReservationService Server
	proto.RegisterReservationServiceServer(gcrpServer, &server{})
	reflection.Register(gcrpServer)

	if err := gcrpServer.Serve(lis); err != nil {
		log.Fatalf("Error %v", err)
	}

}

type server struct{}

func (s *server) MakeReservation(ctx context.Context, req *proto.BookTrip) (*proto.Trip, error) {
	log.Printf("Recieved request to book trip from client %s", req.PassengerName)
	new_trip := proto.Trip{
		PassengerName: req.PassengerName,
		DriverName:    "Marek",
	}
	return &new_trip, nil
}
