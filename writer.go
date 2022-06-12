package persistencesql

import (
	"github.com/asynkron/protoactor-go/actor"
)

type write struct {
	fun func()
}

func newWriter() func(actor.Context) {
	return func(context actor.Context) {
		switch msg := context.Message().(type) {
		case *write:
			go msg.fun()
		}
	}
}
