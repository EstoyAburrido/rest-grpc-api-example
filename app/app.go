package app

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"time"

	fibgrpc "github.com/estoyaburrido/rest-grpc-api-example/app/grpc/controllers/fib-grpc"
	"github.com/estoyaburrido/rest-grpc-api-example/app/rest/controllers/fib"
	"github.com/estoyaburrido/rest-grpc-api-example/app/tools"
	"github.com/sirupsen/logrus"
)

type App struct {
	host      string
	httpPort  string
	grpcPort  string
	fibonacci *tools.Fibonacci
}

func NewApp(
	host string,
	httpPort string,
	grpcPort string,
	fibonacci *tools.Fibonacci,
) *App {
	return &App{
		host:      host,
		httpPort:  httpPort,
		grpcPort:  grpcPort,
		fibonacci: fibonacci,
	}
}

func (app *App) Run() {
	const serviceNum = 2
	var wg sync.WaitGroup

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	wg.Add(serviceNum)

	notifyQuit := make([]chan struct{}, serviceNum)
	for i := range notifyQuit {
		notifyQuit[i] = make(chan struct{}, 1)
	}
	go app.runFibRest(&wg, notifyQuit[0])
	go app.runFibGrpc(&wg, notifyQuit[1])

	<-quit
	for i := range notifyQuit {
		notifyQuit[i] <- struct{}{}
	}

	wg.Wait()
}

func (app *App) runFibRest(wg *sync.WaitGroup, quit chan struct{}) {
	c := fib.NewController(
		app.host, app.httpPort, app.fibonacci,
	)

	errChan := make(chan error, 1)
	go func() {
		if err := c.Run(); err != nil {
			errChan <- err
		}
	}()

	select { // stopping on both an error or an SIGINT
	case err := <-errChan:
		c.Logger().Printf("Fatal error: %v\n", err)
	case <-quit:
		c.Logger().Printf("Shutting down Fib-REST gracefully")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := c.Shutdown(ctx); err != nil {
		c.Logger().Fatal(err)
	}

	wg.Done()
}

func (app *App) runFibGrpc(wg *sync.WaitGroup, quit chan struct{}) {
	grpcLogger := logrus.New().WithField("Service", "GRPC-Fib")

	c := fibgrpc.NewController(
		app.host, app.grpcPort, grpcLogger, app.fibonacci,
	)

	errChan := make(chan error)
	go func() {
		if err := c.Run(); err != nil {
			errChan <- err
		}
	}()

	select { // stopping on both an error or an SIGINT
	case err := <-errChan:
		grpcLogger.Printf("Fatal error: %v\n", err)
	case <-quit:
		grpcLogger.Printf("Shutting down Fib-gRPC gracefully\n")
	}

	c.Shutdown()

	wg.Done()

}
