package transport

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/benkim0414/lott/pkg/endpoint"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

var ErrBadRouting = errors.New("inconsistent mapping between route and handler (programmer error)")

func NewHTTPHandler(endpoints endpoint.Endpoints, logger log.Logger) http.Handler {
	r := mux.NewRouter().PathPrefix("/api/v0/").Subrouter()
	options := []httptransport.ServerOption{
		httptransport.ServerErrorLogger(logger),
		httptransport.ServerErrorEncoder(encodeError),
	}
	r.Methods("GET").Path("/products/{productID}/draws").Handler(httptransport.NewServer(
		endpoints.ListDrawsEndpoint,
		decodeListDrawsRequest,
		encodeResponse,
		options...,
	))
	r.Methods("GET").Path("/products/{productID}/draws/{drawID}").Handler(httptransport.NewServer(
		endpoints.GetDrawEndpoint,
		decodeGetDrawRequest,
		encodeResponse,
		options...,
	))
	return r
}

func decodeListDrawsRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	productID, ok := vars["productID"]
	if !ok {
		return nil, ErrBadRouting
	}
	return endpoint.ListDrawsRequest{ProductID: productID}, nil
}

func decodeGetDrawRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	productID, ok := vars["productID"]
	if !ok {
		return nil, ErrBadRouting
	}
	drawID, ok := vars["drawID"]
	if !ok {
		return nil, ErrBadRouting
	}
	return endpoint.GetDrawRequest{ProductID: productID, DrawID: drawID}, nil
}

type errorer interface {
	error() error
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
