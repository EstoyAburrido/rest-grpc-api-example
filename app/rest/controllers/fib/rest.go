package fib

import (
	"context"
	"fmt"
	"net/http"

	"github.com/estoyaburrido/rest-grpc-api-example/app/rest/models"
	"github.com/estoyaburrido/rest-grpc-api-example/app/tools"
	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Controller struct {
	fibonacci *tools.Fibonacci
	e         *echo.Echo
	host      string
	port      string
}

func NewController(host, port string, fibonacci *tools.Fibonacci) *Controller {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	c := &Controller{
		e:         e,
		fibonacci: fibonacci,
		host:      host,
		port:      port,
	}

	c.e.POST("/getSequence", c.getFibonacci)
	c.e.POST("/getMaxIndex", c.getMaxIndex)

	return c
}

func (controller *Controller) Run() error {
	connectionString := fmt.Sprintf("%v:%v", controller.host, controller.port)
	return controller.e.Start(connectionString)
}

func (controller *Controller) Shutdown(ctx context.Context) error {
	return controller.e.Shutdown(ctx)
}

func (controller *Controller) Logger() echo.Logger {
	return controller.e.Logger
}

// For testing
func (controller *Controller) getMaxIndex(c echo.Context) error {
	res := controller.fibonacci.GetMaxIndex()

	resp := models.FibbonaciIndex{
		MaxIndex: *res,
	}

	return c.JSON(http.StatusOK, resp)
}

func (controller *Controller) getFibonacci(c echo.Context) error {
	req := &models.FibbonaciQuery{}
	if err := c.Bind(req); err != nil {
		return err
	}

	res := controller.fibonacci.GetSequence(req.X, req.Y)

	return c.JSON(http.StatusOK, res)
}
