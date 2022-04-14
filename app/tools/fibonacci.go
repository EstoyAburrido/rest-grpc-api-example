package tools

import (
	log "github.com/sirupsen/logrus"

	"github.com/estoyaburrido/rest-grpc-api-example/app/databases"
)

type Fibonacci struct {
	redis  databases.Redis
	logger *log.Entry
}

func NewFibonacci(redis databases.Redis, logger *log.Entry) *Fibonacci {
	return &Fibonacci{
		redis:  redis,
		logger: logger,
	}
}

func (f *Fibonacci) GetSequence(from, to uint64) []uint64 {
	var res []uint64
	currentMax, err := f.redis.GetMaxFibonacci()

	usingCache := true
	if err != nil {
		f.logger.Println("Can't retrieve currently saved maximum index of the Fibonacci sequence")
		f.logger.Println(err.Error())

		usingCache = false
	} else if *currentMax < to {
		f.logger.Println("The stored sequence length is less than needed, extending")

		cachedValue, err := f.redis.GetSavedFibonacci()
		if err != nil {
			f.logger.Errorln("Can't retrieve currently saved Fibonacci sequence")
			f.logger.Errorln(err.Error())

			usingCache = false
		} else {
			res = f.extend(cachedValue, to)

			err := f.redis.SaveFibonacci(res)
			if err != nil {
				f.logger.Errorln("failed to save new Fibonacci sequence")
				f.logger.Errorln(err.Error())
			}

			return res[from : to+1]
		}
	}

	if usingCache {
		cachedValue, err := f.redis.GetSavedFibonacci()
		if err != nil {
			f.logger.Errorln("Can't retrieve currently saved Fibonacci sequence")
			f.logger.Errorln(err.Error())
		} else {
			res = cachedValue[from : to+1]

			return res
		}
	}

	res = f.recalculate(to)

	err = f.redis.SaveFibonacci(res)
	if err != nil {
		f.logger.Errorln("failed to save new Fibonacci sequence")
		f.logger.Errorln(err.Error())
	}

	return res[from : to+1]

}

func (f *Fibonacci) GetMaxIndex() *uint64 {
	res, err := f.redis.GetMaxFibonacci()
	if err != nil {
		return nil
	}

	return res
}

func (f *Fibonacci) recalculate(to uint64) []uint64 {
	return f.extend(nil, to)
}

func (f *Fibonacci) extend(from []uint64, to uint64) []uint64 {
	res := make([]uint64, 0)

	var startIndex uint64

	if from != nil {
		maxIndex := len(from) - 1
		if maxIndex < 0 {
			maxIndex = 0
		}
		startIndex = uint64(maxIndex)
		res = from
	} else {
		res = append(res, 0, 1, 1)
		startIndex = 2
	}

	for i := startIndex; i < to; i++ {
		res = append(res, res[i]+res[i-1])
	}

	return res
}
