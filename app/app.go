package app

import (
	"fmt"
	"sync"

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
	var wg sync.WaitGroup

	wg.Add(2)
	go app.runRest(&wg)
	go app.runGrpc(&wg)
	wg.Wait()
}

func (app *App) runRest(wg *sync.WaitGroup) {
	c := fib.NewController(
		app.host, app.httpPort, app.fibonacci,
	)

	_ = fmt.Errorf(c.Run().Error())

	wg.Done()
}

func (app *App) runGrpc(wg *sync.WaitGroup) {
	grpcLogger := logrus.New().WithField("Service", "GRPC")

	c := fibgrpc.NewController(
		app.host, app.grpcPort, grpcLogger, app.fibonacci,
	)

	_ = fmt.Errorf(c.Run().Error())

	wg.Done()

}
