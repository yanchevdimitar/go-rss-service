package queue

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

type Queue interface {
	GracefullyStop() chan bool
}

type Processor interface {
	Process()
}

type DefaultQueue struct {
}

func NewDefaultQueue() Queue {
	return DefaultQueue{}
}

func (q DefaultQueue) GracefullyStop() chan bool {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT)
	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		done <- true
	}()

	return done
}
