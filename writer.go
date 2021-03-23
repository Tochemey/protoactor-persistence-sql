package persistence

import (
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
)

type write struct {
	fun func()
}

func newWriter(rate time.Duration) func(actor.Context) {
	return func(context actor.Context) {
		switch msg := context.Message().(type) {
		case *write:
			go msg.fun()
			//	time.Sleep(rate)
		}
	}
}
