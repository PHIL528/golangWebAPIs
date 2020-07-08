package rest

import (
	"bytes"
	"encoding/json"
	//"errors"
	//"github.com/marchmiel/proto-playground/Config"
	"github.com/marchmiel/proto-playground/client/model"
	//"github.com/marchmiel/proto-playground/proto"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"time"
)

type restTripBooker struct {
	httpCli *http.Client
}

//cli *http.Client
func NewTripBooker() (model.TripBooker, error) {
	return &restTripBooker{httpCli: &http.Client{}}, nil
}
func (r *restTripBooker) BookTrip(mod *model.BookTripRequest) (*model.TripBookedResponse, error) {
	bookTripRequest, err := CreateBuffer(mod)
	if err != nil {
		return nil, errors.Wrap(err, "Could not create JSON buffer")
	}
	req, err := http.NewRequest("POST", "http://localhost:3003/", bookTripRequest)
	fmt.Println(req)
	req.Header.Set("Content-Type", "application/json")
	r.httpCli.Timeout = 10 * time.Second
	resp, err := r.httpCli.Do(req)
	//defer resp.Body.Close()
	if err != nil {
		return nil, errors.Wrap(err, "Could not fetch http request")
	}
	var tripBookedResponse model.TripBookedResponse
	err = json.NewDecoder(resp.Body).Decode(&tripBookedResponse)
	fmt.Println(tripBookedResponse)
	if err != nil {
		return nil, errors.Wrap(err, "Could not decode JSON")
	}
	return &tripBookedResponse, nil
}
func CreateBuffer(mod *model.BookTripRequest) (*bytes.Buffer, error) {
	jsonbytes, err := json.Marshal(mod)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(jsonbytes), nil
}
