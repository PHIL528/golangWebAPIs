package servrest

import (
	//"bytes"
	"encoding/json"
	//	"fmt"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/marchmiel/proto-playground/client/model"
	//"github.com/marchmiel/proto-playground/server/wrapper"
	"net/http"
)

type restDataType struct {
	req  *http.Request
	resp http.ResponseWriter
	err  chan error
}

func NewRestData(rw http.ResponseWriter, r *http.Request) *restDataType {
	return &restDataType{req: r, resp: rw, err: make(chan error)}
}

func (r *restDataType) Unload(btr *model.BookTripRequest) error {
	return json.NewDecoder(r.req.Body).Decode(btr)
}
func (r *restDataType) Load(tbr *model.TripBookedResponse) error {
	r.resp.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(r.resp).Encode(tbr)
	r.SendBack(err)
	return err
}
func (r *restDataType) GetResponse() interface{} {
	return nil
}
func (r *restDataType) CorrelationID() string {
	return watermill.NewUUID()
}
func (r *restDataType) SendBack(err error) {
	r.err <- err
}

//jsonbytes, err := json.Marshal(tbr)
//r.resp.Body = bytes.NewBuffer(jsonbytes)
//err := errors.New("mock")
