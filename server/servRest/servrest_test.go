package servrest

import (
	"bytes"
	"encoding/json"
	"github.com/marchmiel/proto-playground/client/model"
	"github.com/marchmiel/proto-playground/server/wrapper"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var testPassenger string = "RestServerPassenger404"

func TestServRest(t *testing.T) {
	collection := make(chan wrapper.ClientDataType, 100)
	rHandler := &restHandler{collection}
	testServer := httptest.NewServer(rHandler)
	defer func() { testServer.Close() }()
	time.Sleep(time.Second)

	btrInput := &model.BookTripRequest{
		PassengerName: testPassenger,
	}
	jsonbytes, _ := json.Marshal(btrInput)
	testCli := &http.Client{}
	var resp *http.Response
	respFinished := make(chan int)
	go func() {
		resp, _ = testCli.Post(testServer.URL, "application/json", bytes.NewBuffer(jsonbytes))
		respFinished <- 747
	}()

	var invoke wrapper.ClientDataType = nil
	select {
	case <-time.After(time.Second * 5):
		t.Error("Timeout in server")
	case invk := <-collection:
		invoke = invk
	}
	var btrInvoke model.BookTripRequest
	invoke.Unload(&btrInvoke)
	assert.Equal(t, btrInput, &btrInvoke, " compaing BTRequest sent to BTRequest recieved by collection channel")
	tbrStub := &model.TripBookedResponse{
		PassengerName: testPassenger,
		DriverName:    "TestDriver404",
	}
	invoke.Load(tbrStub)
	select {
	case <-time.After(time.Second * 5):
		t.Error("Timeout waiting for response")
	case <-respFinished:
		var tripBookedResponse model.TripBookedResponse
		json.NewDecoder(resp.Body).Decode(&tripBookedResponse)
		assert.Equal(t, tbrStub, &tripBookedResponse, " comparing TBResp given to server to that recieved by client")
	}

}
