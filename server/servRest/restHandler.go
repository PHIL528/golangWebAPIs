package servrest

import (
	"github.com/marchmiel/proto-playground/server/wrapper"
	"github.com/pkg/errors"
	"log"
	"net/http"
)

type restHandler struct {
	post chan<- wrapper.ClientDataType
}

func NewServePortChan(PORT string, Post chan<- wrapper.ClientDataType, sendback chan<- error) {
	rHandler := &restHandler{Post}
	mux := http.NewServeMux()
	mux.Handle("/", rHandler)
	err := http.ListenAndServe(":3003", mux)
	if err != nil {
		sendback <- errors.Wrap(err, "Could not listen and serve")
	}
}

func (r *restHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	clientPkg := NewRestData(rw, req)
	r.post <- clientPkg
	err := <-clientPkg.err
	if err != nil {
		log.Printf(" error in ServeHTTP method")
	}
}
