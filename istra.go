package istra

import (
	"errors"
)

var (
	ErrCreatingChannel  = errors.New("error creating channel")
	ErrConsumingChannel = errors.New("error consuming queue")
)

type Connection interface {
	Channel() (Channel, error)
	NotifyOnClose() chan error
}

type Messenger interface {
	Msg() []byte
}

type Channel interface {
	Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args interface{}) (<-chan Messenger, error)
}

type QueueConf struct {
	Name      string
	AutoAck   bool
	Exclusive bool
	NoLocal   bool
	NoWait    bool
}

func ProcessQueue(amqp Connection, conf QueueConf, f func([]byte)) {
	ch, err := amqp.Channel()
	if err != nil {
		panic(ErrCreatingChannel)
	}

	messenger, err := ch.Consume(conf.Name, "", conf.AutoAck, conf.Exclusive, conf.NoLocal, conf.NoWait, nil)
	if err != nil {
		panic(ErrConsumingChannel)
	}

	for {
		select {
		case m := <-messenger:
			f(m.Msg())
		case _ = <-amqp.NotifyOnClose():
			return
		}
	}
}
