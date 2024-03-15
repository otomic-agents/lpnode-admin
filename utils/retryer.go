package utils

import (
	"log"
	"time"

	"github.com/pkg/errors"
)

type RepetOption struct {
	Interval int64
	MaxCount int64
}
type Retryer struct {
	Option *RepetOption
}

func RetryerNew() *Retryer {
	return &Retryer{}
}
func (rt *Retryer) SetOption(option *RepetOption) *Retryer {
	rt.Option = option
	return rt
}

func (rt *Retryer) Repet(doFun func() error) (err error) {
	var i = 0
	var lastErr error
	for {
		i++
		log.Printf("attempt no %d", i)
		if i >= int(rt.Option.MaxCount) {
			err = errors.WithMessage(GetNoEmptyError(lastErr), "exceeded max tries, final error was:")
			return
		}
		execErr := doFun()
		if execErr != nil {
			lastErr = execErr
			time.Sleep(time.Millisecond * time.Duration(rt.Option.Interval))
		} else {
			return
		}
	}
}
