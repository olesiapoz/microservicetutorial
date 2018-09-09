package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/olesiapoz/microservicetutorial"
	"golang.org/x/net/context"
)

func main() {

	ctx := context.Background()

	errChan := make(chan error)

	var svc microservicetutorial.Service

	svc = microservicetutorial.LoremService{}

	endpoint := microservicetutorial.Endpoints{

		LoremEndpoint: microservicetutorial.MakeLoremEndpoint(svc),
	}

	// Logging domain.

	var logger log.Logger

	{

		logger = log.NewLogfmtLogger(os.Stderr)

		logger = log.With(logger, "ts", log.DefaultTimestampUTC)

		logger = log.With(logger, "caller", log.DefaultCaller)

	}

	r := microservicetutorial.MakeHttpHandler(ctx, endpoint, logger)

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
