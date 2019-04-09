package service

import (
	"context"

	"github.com/go-kit/kit/log"
)

type Middleware func(Service) Service

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return loggingMiddleware{logger, next}
	}
}

type loggingMiddleware struct {
	logger log.Logger
	next   Service
}

func (mw loggingMiddleware) ListDraws(ctx context.Context, productID string) (draws []*Draw, err error) {
	defer func() {
		mw.logger.Log("method", "ListDraws", "productID", productID, "err", err)
	}()
	return mw.next.ListDraws(ctx, productID)
}

func (mw loggingMiddleware) GetDraw(ctx context.Context, productID string, drawID string) (draw *Draw, err error) {
	defer func() {
		mw.logger.Log("method", "GetDraw", "productID", productID, "drawID", drawID, "err", err)
	}()
	return mw.next.GetDraw(ctx, productID, drawID)
}
