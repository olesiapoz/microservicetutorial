package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ru-rocker/gokit-playground/lorem"

	"github.com/go-kit/kit/log"
	"golang.org/x/net/context"
)

func main() {
	ctx := context.Background()
	errChan := make(chan error)

	var svc lorem.Service
	svc = lorem.LoremService{}
	endpoint := lorem.Endpoints{
		LoremEndpoint: lorem.MakeLoremEndpoint(svc),
	}

	// Logging domain.
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.NewContext(logger).With("ts", log.DefaultTimestampUTC)
		logger = log.NewContext(logger).With("caller", log.DefaultCaller)
	}

	r := lorem.MakeHttpHandler(ctx, endpoint, logger)

	// HTTP transport
	go func() {
		fmt.Println("Starting server at port 8080")
		handler := r
		errChan <- http.ListenAndServe(":8080", handler)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()
	fmt.Println(<-errChan)
}
