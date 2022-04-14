package fibgrpc

import (
	"context"
	"testing"

	pb "github.com/estoyaburrido/rest-grpc-api-example/app/grpc/proto"
	"github.com/estoyaburrido/rest-grpc-api-example/app/rest/models"
	"github.com/estoyaburrido/rest-grpc-api-example/app/tools"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

type redisMock struct {
	storedSequence []uint64
}

func (r *redisMock) GetMaxFibonacci() (*uint64, error) {
	res, err := r.GetSavedFibonacci()
	if err != nil {
		return nil, err
	}

	length := len(res) - 1
	if length < 0 {
		length = 0
	}

	retVal := uint64(length)
	return &retVal, nil
}

func (r *redisMock) GetSavedFibonacci() ([]uint64, error) {
	return r.storedSequence, nil
}

func (r *redisMock) SaveFibonacci(sequence []uint64) error {
	r.storedSequence = sequence

	return nil
}

var reference = []uint64{0, 1, 1, 2, 3, 5, 8, 13, 21, 34, 55, 89, 144, 233}

var firstReq = models.FibbonaciQuery{
	X: 0,
	Y: 5,
}
var secondReq = models.FibbonaciQuery{
	X: 2,
	Y: 13,
}

func TestGrpc(t *testing.T) {

	loggerFib := logrus.New().WithField("Service", "Fibonacci")
	fib := tools.NewFibonacci(&redisMock{}, loggerFib)
	loggerGrpc := logrus.New().WithField("Service", "gRPC")
	controller := &GrpcController{
		logger:    loggerGrpc,
		fibonacci: fib,
	}

	req := &pb.GetSequenceRequest{
		X: firstReq.X,
		Y: firstReq.Y,
	}

	// Testing from 0 idx
	res, err := controller.Get(context.Background(), req)
	if assert.NoError(t, err, "error at 1st grpc test") {
		assert.ElementsMatch(t, reference[:firstReq.Y+1], res.Res, "elements don't match at 1st grpc test")
	}

	// Testing from non-zero index and expanding the cached array
	req = &pb.GetSequenceRequest{
		X: secondReq.X,
		Y: secondReq.Y,
	}
	res, err = controller.Get(context.Background(), req)
	if assert.NoError(t, err, "error at 2nd grpc test") {
		assert.ElementsMatch(t, reference[secondReq.X:secondReq.Y+1], res.Res, "elements don't match at 2nd grpc test")
	}

}
