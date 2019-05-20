package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"cloud.google.com/go/firestore"
	"github.com/benkim0414/lott/pkg/endpoint"
	"github.com/benkim0414/lott/pkg/service"
	"github.com/benkim0414/lott/pkg/transport"
	"github.com/go-kit/kit/log"
	"github.com/gorilla/handlers"
)

func main() {
	var (
		httpAddr = flag.String("http.addr", ":8080", "HTTP listen address")
	)
	flag.Parse()

	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)

	ctx := context.Background()
	projectID := os.Getenv("PROJECT_ID")
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		logger.Log("failed to create Firestore client: %v", err)
	}
	defer client.Close()

	var (
		svc         = service.New(client, logger)
		endpoints   = endpoint.New(svc, logger)
		httpHandler = transport.NewHTTPHandler(endpoints, logger)
	)

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		logger.Log("transport", "HTTP", "addr", *httpAddr)
		errs <- http.ListenAndServe(*httpAddr, handlers.CORS()(httpHandler))
	}()
	logger.Log("exit", <-errs)
}
