package grpc

import (
	//	"context"
	//"fmt"
	"github.com/marchmiel/proto-playground/client/mocks"
	"github.com/marchmiel/proto-playground/client/model"
	"github.com/marchmiel/proto-playground/proto"
	"github.com/stretchr/testify/assert"
	"testing"
)

var testPassenger string = "GrpcTestPassenger"
var testDriver string = "GrpcDriver202"

//Testcase information - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
var TestBookTripRequest *model.BookTripRequest = &model.BookTripRequest{
	PassengerName: testPassenger,
}
var ExpectedInvocatedProto *proto.BookTrip = &proto.BookTrip{
	PassengerName: TestBookTripRequest.PassengerName,
}
var MockReturnedProto *proto.TripBooked = &proto.TripBooked{
	Trip: &proto.Trip{
		PassengerName: ExpectedInvocatedProto.PassengerName,
		DriverName:    testDriver,
	},
}
var ExpectedTrip *model.TripBookedResponse = &model.TripBookedResponse{
	PassengerName: MockReturnedProto.Trip.PassengerName,
	DriverName:    MockReturnedProto.Trip.DriverName,
}

func TestAsMockGrpc(t *testing.T) {
	mockRes := &mocks.FakeReservationServiceClient{}
	mockRes.MakeReservationReturns(MockReturnedProto, nil)

	tripBooker, _ := CustomTripBooker(mockRes)
	tripBooked, _ := tripBooker.BookTrip(TestBookTripRequest)

	assert.Equal(t, ExpectedInvocatedProto /*&proto.BookTrip{PassengerName: testPassenger}*/, mockRes.Invocations()["MakeReservation"][0][1].(*proto.BookTrip), " comparing BookTrip proto sent through mock MakeReservation() service to expected value")
	assert.Equal(t, ExpectedTrip, tripBooked, " comparing TripBooked returned by BookTrip() in mock gRPC client to expected value")
}
