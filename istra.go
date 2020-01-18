package istra

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

var (
	ErrCreatingChannel  = errors.New("error creating channel")
	ErrConsumingChannel = errors.New("error consuming queue")
	UnknownBinding      = errors.New("unknown binding; allowed: 'Declare', 'Bind', 'UnBind', 'DeclareBind'")
)

type connection interface {
	consumeChanneler
	notifyOnClose() chan *amqp.Error
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
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("declaring queue : '%s' failed", b.(Declare).Name))
			}
		case Bind:
			err = ch.bind(b.(Bind))
			if err != nil {
				return wrap(err, "binding", b.(Bind).Name)
			}
		case UnBind:
			err = ch.unbind(b.(UnBind))
			if err != nil {
				return wrap(err, "unbinding", b.(UnBind).Queue)
			}
		case DeclareBind:
			err = ch.declare(b.(DeclareBind).Declare)
			if err != nil {
				return wrap(err, "declaring", b.(DeclareBind).Declare.Name)
			}
			err = ch.bind(b.(DeclareBind).Bind)
			if err != nil {
				return wrap(err, "binding", b.(DeclareBind).Bind.Name)
			}
		default:
			return UnknownBinding
		}
	}
	return nil
}

func wrap(err error, method, queue string) error {
	return errors.Wrap(err, fmt.Sprintf("%s queue : '%s' failed", method, queue))
}
