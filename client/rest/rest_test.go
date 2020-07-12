package rest

import (
	"bytes"
	"encoding/json"
	"github.com/marchmiel/proto-playground/client/model"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

var testPassengerName string = "restTestPassenger303"

func TestRestWithMockServer(t *testing.T) {
	invocation := make(chan *bytes.Buffer, 2)
	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		bys, _ := ioutil.ReadAll(req.Body)
		invocation <- bytes.NewBuffer(bys)
		//Response is not function of input
		json.NewEncoder(res).Encode(StubResponse)
	}))
	defer func() { testServer.Close() }()
	url = testServer.URL
	restTripBooker, _ := NewTripBooker()
	tripBooked, _ := restTripBooker.BookTrip(testBookTripRequest)
	assert.Equal(t, expectedInvocation(), <-invocation, " checking request made to server")
	assert.Equal(t, tripBooked, StubResponse, " comparing StubResponse to encoded/decoded StubResponse")
}

var testBookTripRequest *model.BookTripRequest = &model.BookTripRequest{
	PassengerName: testPassengerName,
}
var expectedInvocation = func() *bytes.Buffer {
	jsonbytes, _ := json.Marshal(testBookTripRequest)
	return bytes.NewBuffer(jsonbytes)
}
var StubResponse = &model.TripBookedResponse{
	PassengerName: testBookTripRequest.PassengerName,
	DriverName:    "restTestDriver30303",
}

/*
type HTTPClientInterface interface {
	Do(req *http.Request) (*http.Response, error)
}

type restTripBooker struct {
	httpCli HTTPClientInterface
}

func NewTripBooker() (model.TripBooker, error) {
	return &restTripBooker{httpCli: &http.Client{}}, nil
}
func CustomTripBooker(httpCliIntf HTTPClientInterface) (model.TripBooker, error) {
	return &restTripBooker{httpCli: httpCliInft}, nil
}
*/

/*
	jsonbytes, _ := json.Marshal(testBookTripRequest)
	buff := bytes.NewBuffer(jsonbytes)
	buf := buff.Bytes()
	r := bytes.NewReader(buf)
	return ioutil.NopCloser(r)
*/
