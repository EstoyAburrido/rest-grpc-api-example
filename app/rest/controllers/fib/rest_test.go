package fib

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/estoyaburrido/rest-grpc-api-example/app/rest/models"
	"github.com/estoyaburrido/rest-grpc-api-example/app/tools"
	echo "github.com/labstack/echo/v4"
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

func TestRest(t *testing.T) {
	logger := logrus.New().WithField("Service", "Fibonacci")
	fib := tools.NewFibonacci(&redisMock{}, logger)

	e := echo.New()

	h := &Controller{
		fibonacci: fib,
		e:         e,
	}

	// Testing from 0 idx
	reqJSON, _ := json.Marshal(firstReq)

	req := httptest.NewRequest(http.MethodPost, "/getSequence", bytes.NewReader(reqJSON))
	req.Header["Content-Type"] = []string{"application/json"}

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Response().Header()["Accept"] = []string{"application/json"}

	if assert.NoError(t, h.getFibonacci(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		var seqResp []uint64
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &seqResp))
		assert.ElementsMatch(t, reference[:firstReq.Y+1], seqResp)
	}

	// Testing from non-zero index and expanding the cached array
	reqJSON, _ = json.Marshal(secondReq)
	req = httptest.NewRequest(http.MethodPost, "/getSequence", bytes.NewReader(reqJSON))
	req.Header["Content-Type"] = []string{"application/json"}

	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.Response().Header()["Accept"] = []string{"application/json"}

	if assert.NoError(t, h.getFibonacci(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		var seqResp []uint64
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &seqResp))
		assert.ElementsMatch(t, reference[secondReq.X:secondReq.Y+1], seqResp)
	}
}
