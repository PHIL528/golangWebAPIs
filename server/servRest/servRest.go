package servrest

import (
	//"bytes"
	"encoding/json"
	//	"fmt"
	//	"github.com/marchmiel/proto-playground/client/model"
	"github.com/marchmiel/proto-playground/server/wrapper"
	//"io"
	//	"log"
	//"github.com/pkg/errors"
	"net/http"
)

type restData struct {
	req  *http.Request
	resp *http.ResponseWriter
}

func NewRestData(rw *http.ResponseWriter, r *http.Request) *restData {
	return &restData{req: r, resp: rw}
}

func (r *restData) Unload(btr *wrapper.BookTripRequest) error {
	return json.NewDecoder(r.req.Body).Decode(btr)
}
func (r *restData) Load(tbr *wrapper.TripBookedResponse) error {
	//jsonbytes, err := json.Marshal(tbr)

	//r.resp.Body = bytes.NewBuffer(jsonbytes)
	//err := errors.New("mock")
	return nil
}
func (g *restData) Ret() interface{} {
	return nil
}
