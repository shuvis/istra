package istra

import (
	"errors"
)

var (
	ErrCreatingChannel  = errors.New("error creating channel")
	ErrConsumingChannel = errors.New("error consuming queue")
)

type connection interface {
	channel() (channel, error)
	notifyOnClose() chan error
}

type messenger interface {
	msg() []byte
}

type channel interface {
	consume(queue string, autoAck, exclusive, noLocal, noWait bool) (<-chan messenger, error)
}

type QueueConf struct {
	Name      string
	AutoAck   bool
	Exclusive bool
	NoLocal   bool
	NoWait    bool
}

func processQueue(amqp connection, conf QueueConf, f func([]byte)) {
	ch, err := amqp.channel()
	if err != nil {
		panic(ErrCreatingChannel)
	}

	messenger, err := ch.consume(conf.Name, conf.AutoAck, conf.Exclusive, conf.NoLocal, conf.NoWait)
	if err != nil {
		panic(ErrConsumingChannel)
	}

	for {
		select {
		case m := <-messenger:
			f(m.msg())
		case _ = <-amqp.notifyOnClose():
			return
		}
	}
}
