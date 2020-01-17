package istra

import (
	"errors"
	"github.com/streadway/amqp"
)

var (
	ErrCreatingChannel  = errors.New("error creating channel")
	ErrConsumingChannel = errors.New("error consuming queue")
)

type connection interface {
	channel() (channel, error)
	notifyOnClose() chan error
}

type channel interface {
	consume(queue string, autoAck, exclusive, noLocal, noWait bool) (<-chan amqp.Delivery, error)
}

func processQueue(amqp connection, conf QueueConf, f func(amqp.Delivery)) {
	ch, err := amqp.channel()
	if err != nil {
		panic(ErrCreatingChannel)
	}

	deliveries, err := ch.consume(conf.Name, conf.AutoAck, conf.Exclusive, conf.NoLocal, conf.NoWait)
	if err != nil {
		panic(ErrConsumingChannel)
	}

	for {
		select {
		case d := <-deliveries:
			f(d)
		case _ = <-amqp.notifyOnClose():
			return
		}
	}
}
