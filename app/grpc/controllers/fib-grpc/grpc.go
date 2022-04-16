package fibgrpc

import (
	"context"
	"errors"
	"fmt"
	"net"

	pb "github.com/estoyaburrido/rest-grpc-api-example/app/grpc/proto"
	"github.com/estoyaburrido/rest-grpc-api-example/app/tools"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type GrpcController struct {
	fibonacci  *tools.Fibonacci
	host       string
	port       string
	grpcServer *grpc.Server
	pb.UnimplementedFibonacciServiceServer
	logger *logrus.Entry
}

func NewController(host, port string, logger *logrus.Entry, fibonacci *tools.Fibonacci) *GrpcController {
	s := grpc.NewServer()
	controller := &GrpcController{
		logger:     logger,
		host:       host,
		port:       port,
		grpcServer: s,
		fibonacci:  fibonacci,
	}
	pb.RegisterFibonacciServiceServer(s, controller)

	return controller
}

func (c *GrpcController) Shutdown() {
	c.grpcServer.GracefulStop()
}

func (c *GrpcController) Run() error {
	lis, err := net.Listen("tcp", fmt.Sprintf("%v:%v", c.host, c.port))
	if err != nil {
		c.logger.Fatalf("gprc failed to listen: %v", err)
		return err
	}

	c.logger.Printf("server listening at %v", lis.Addr())

	if err := c.grpcServer.Serve(lis); err != nil {
		c.logger.Fatalf("failed to serve: %v", err)
		return err
	}

	return errors.New("lis.Accept failed")
}

func (c *GrpcController) Get(ctx context.Context, req *pb.GetSequenceRequest) (*pb.GetSequenceResponse, error) {
	seq := c.fibonacci.GetSequence(req.X, req.Y)

	res := &pb.GetSequenceResponse{
		Res: seq,
	}

	return res, nil
}
