package tools

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFib(t *testing.T) {
	assert := assert.New(t)

	reference := []uint64{0, 1, 1, 2, 3, 5, 8, 13, 21, 34, 55, 89, 144, 233}

	dummyObj := &Fibonacci{}
	res := dummyObj.recalculate(10)

	for i := range res {
		assert.Equal(res[i], reference[i], "The element value is incorrect")
	}

	res = dummyObj.extend(res, 13)

	for i := range res {
		assert.Equal(res[i], reference[i], "The element value is incorrect")
	}
}
