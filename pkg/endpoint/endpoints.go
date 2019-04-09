package endpoint

import (
	"context"

	"github.com/benkim0414/lott/pkg/service"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
)

type Endpoints struct {
	ListDrawsEndpoint endpoint.Endpoint
	GetDrawEndpoint   endpoint.Endpoint
}

func New(svc service.Service, logger log.Logger) Endpoints {
	listDrawsEndpoint := MakeListDrawsEndpoint(svc)
	listDrawsEndpoint = LoggingMiddleware(log.With(logger, "method", "ListDraws"))(listDrawsEndpoint)

	getDrawEndpoint := MakeGetDrawEndpoint(svc)
	getDrawEndpoint = LoggingMiddleware(log.With(logger, "method", "GetDraw"))(getDrawEndpoint)

	return Endpoints{
		ListDrawsEndpoint: listDrawsEndpoint,
		GetDrawEndpoint:   getDrawEndpoint,
	}
}

func MakeListDrawsEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(ListDrawsRequest)
		d, e := svc.ListDraws(ctx, req.ProductID)
		return ListDrawsResponse{Draws: d, Err: e}, nil
	}
}

func MakeGetDrawEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetDrawRequest)
		d, e := svc.GetDraw(ctx, req.ProductID, req.DrawID)
		return GetDrawResponse{Draw: d, Err: e}, nil
	}
}

type ListDrawsRequest struct {
	ProductID string `json:"proudctId"`
}

type ListDrawsResponse struct {
	Draws []*service.Draw `json:"draws"`
	Err   error           `json:"err,omitempty"`
}

func (r ListDrawsResponse) error() error { return r.Err }

type GetDrawRequest struct {
	ProductID string `json:"productId"`
	DrawID    string `json:"drawId"`
}

type GetDrawResponse struct {
	Draw *service.Draw `json:"draw"`
	Err  error         `json:"err,omitempty"`
}

func (r GetDrawResponse) error() error { return r.Err }
