package istra

import (
	"errors"
	"github.com/streadway/amqp"
)

var (
	ErrCreatingChannel  = errors.New("error creating consumer")
	ErrConsumingChannel = errors.New("error consuming queue")
	UnknownBinding      = errors.New("unknown binding; allowed: 'Declare', 'Bind', 'UnBind', 'DeclareBind'")
)

type connection interface {
	consumeChanneler
	notifyOnClose() chan error
}

type consumeChanneler interface {
	channel() (consumer, error)
}

type consumer interface {
	consume(queue string, autoAck, exclusive, noLocal, noWait bool) (<-chan amqp.Delivery, error)
}

func consumeQueue(conn connection, conf QueueConf, f func(amqp.Delivery)) {
	ch, err := conn.channel()
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
		case _ = <-conn.notifyOnClose():
			return
		}
	}
}

type bindChanneler interface {
	channel() (binder, error)
}

type binder interface {
	closer
	declare(d Declare) error
	bind(b Bind) error
	unbind(u UnBind) error
}

type closer interface {
	close()
}

func bindQueues(binder bindChanneler, bindings Bindings) error {
	ch, err := binder.channel()
	if err != nil {
		return ErrCreatingChannel
	}
	defer ch.close()

	for _, b := range bindings {
		switch b.(type) {
		case Declare:
			err = ch.declare(b.(Declare))
		case Bind:
			err = ch.bind(b.(Bind))
		case UnBind:
			err = ch.unbind(b.(UnBind))
		case DeclareBind:
			err = ch.declare(b.(DeclareBind).Declare)
			if err != nil {
				return err
			}
			err = ch.bind(b.(DeclareBind).Bind)
		default:
			return UnknownBinding
		}
		if err != nil {
			return err
		}
	}
	return nil
}
