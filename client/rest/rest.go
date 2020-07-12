package rest

import (
	"bytes"
	"encoding/json"
	//"errors"
	//"github.com/marchmiel/proto-playground/Config"
	"github.com/marchmiel/proto-playground/client/model"
	//"github.com/marchmiel/proto-playground/proto"
	//"fmt"
	"github.com/pkg/errors"
	"net/http"
	"time"
)

type restTripBooker struct {
	httpCli *http.Client
}

func NewTripBooker() (model.TripBooker, error) {
	return &restTripBooker{httpCli: &http.Client{}}, nil
}
func CustomRestTripBooker(h *http.Client) (model.TripBooker, error) {
	return &restTripBooker{httpCli: h}, nil
}

var url string = "http://localhost:3003/"

func (r *restTripBooker) BookTrip(mod *model.BookTripRequest) (*model.TripBookedResponse, error) {
	jsonbytes, err := json.Marshal(mod)
	if err != nil {
		return nil, errors.Wrap(err, " could not marshal JSON")
	}
	r.httpCli.Timeout = 10 * time.Second
	resp, err := r.httpCli.Post(url, "application/json", bytes.NewBuffer(jsonbytes))
	if err != nil {
		return nil, errors.Wrap(err, "Client could not do http post")
	}
	defer resp.Body.Close()
	var tripBookedResponse model.TripBookedResponse
	err = json.NewDecoder(resp.Body).Decode(&tripBookedResponse)
	if err != nil {
		return nil, errors.Wrap(err, "Could not decode JSON")
	}
	return &tripBookedResponse, nil
}

/* I added this but later removed it becuase it is not enough content to wrap inside one function
func CreateBuffer(mod *model.BookTripRequest) (*bytes.Buffer, error) {
	jsonbytes, err := json.Marshal(mod)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(jsonbytes), nil
}*/
