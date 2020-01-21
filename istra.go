package istra

import (
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

const (
	channelError = "error creating channel"
	queueError   = "error consuming testQueue"
	bindingError = "binding failed"
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
		panic(errors.Wrap(err, channelError))
	}

	deliveries, err := ch.consume(conf.Name, conf.AutoAck, conf.Exclusive, conf.NoLocal, conf.NoWait)
	if err != nil {
		panic(errors.Wrap(err, queueError))
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

type operationChanneler interface {
	channel() (operator, error)
}

type operator interface {
	closer
	queue(d QueueDeclare) error
	exchange(ed ExchangeDeclare) error
	bind(b Bind) error
	unbind(u UnBind) error
}

type closer interface {
	close()
}

type Operations []Operation

type Operation interface {
	apply(operator) error
}

func processOperations(operator operationChanneler, operations Operations) error {
	ch, err := operator.channel()
	if err != nil {
		return errors.Wrap(err, channelError)
	}
	defer ch.close()

	for _, b := range operations {
		err := b.apply(ch)
		if err != nil {
			return errors.Wrap(err, bindingError)
		}
	}
	return nil
}
