// Code generated by counterfeiter. DO NOT EDIT.
package mocks

import (
	"context"
	"sync"

	"github.com/marchmiel/proto-playground/proto"
	"google.golang.org/grpc"
)

type FakeReservationServiceClient struct {
	MakeReservationStub        func(context.Context, *proto.BookTrip, ...grpc.CallOption) (*proto.TripBooked, error)
	makeReservationMutex       sync.RWMutex
	makeReservationArgsForCall []struct {
		arg1 context.Context
		arg2 *proto.BookTrip
		arg3 []grpc.CallOption
	}
	makeReservationReturns struct {
		result1 *proto.TripBooked
		result2 error
	}
	makeReservationReturnsOnCall map[int]struct {
		result1 *proto.TripBooked
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeReservationServiceClient) MakeReservation(arg1 context.Context, arg2 *proto.BookTrip, arg3 ...grpc.CallOption) (*proto.TripBooked, error) {
	fake.makeReservationMutex.Lock()
	ret, specificReturn := fake.makeReservationReturnsOnCall[len(fake.makeReservationArgsForCall)]
	fake.makeReservationArgsForCall = append(fake.makeReservationArgsForCall, struct {
		arg1 context.Context
		arg2 *proto.BookTrip
		arg3 []grpc.CallOption
	}{arg1, arg2, arg3})
	fake.recordInvocation("MakeReservation", []interface{}{arg1, arg2, arg3})
	fake.makeReservationMutex.Unlock()
	if fake.MakeReservationStub != nil {
		return fake.MakeReservationStub(arg1, arg2, arg3...)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.makeReservationReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeReservationServiceClient) MakeReservationCallCount() int {
	fake.makeReservationMutex.RLock()
	defer fake.makeReservationMutex.RUnlock()
	return len(fake.makeReservationArgsForCall)
}

func (fake *FakeReservationServiceClient) MakeReservationCalls(stub func(context.Context, *proto.BookTrip, ...grpc.CallOption) (*proto.TripBooked, error)) {
	fake.makeReservationMutex.Lock()
	defer fake.makeReservationMutex.Unlock()
	fake.MakeReservationStub = stub
}

func (fake *FakeReservationServiceClient) MakeReservationArgsForCall(i int) (context.Context, *proto.BookTrip, []grpc.CallOption) {
	fake.makeReservationMutex.RLock()
	defer fake.makeReservationMutex.RUnlock()
	argsForCall := fake.makeReservationArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeReservationServiceClient) MakeReservationReturns(result1 *proto.TripBooked, result2 error) {
	fake.makeReservationMutex.Lock()
	defer fake.makeReservationMutex.Unlock()
	fake.MakeReservationStub = nil
	fake.makeReservationReturns = struct {
		result1 *proto.TripBooked
		result2 error
	}{result1, result2}
}

func (fake *FakeReservationServiceClient) MakeReservationReturnsOnCall(i int, result1 *proto.TripBooked, result2 error) {
	fake.makeReservationMutex.Lock()
	defer fake.makeReservationMutex.Unlock()
	fake.MakeReservationStub = nil
	if fake.makeReservationReturnsOnCall == nil {
		fake.makeReservationReturnsOnCall = make(map[int]struct {
			result1 *proto.TripBooked
			result2 error
		})
	}
	fake.makeReservationReturnsOnCall[i] = struct {
		result1 *proto.TripBooked
		result2 error
	}{result1, result2}
}

func (fake *FakeReservationServiceClient) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.makeReservationMutex.RLock()
	defer fake.makeReservationMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeReservationServiceClient) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ proto.ReservationServiceClient = new(FakeReservationServiceClient)
