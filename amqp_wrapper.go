package istra

import "github.com/streadway/amqp"

func ConsumeQueue(conn *amqp.Connection, conf QueueConf, f func(amqp.Delivery)) {
	consumeQueue(&connectionWrapper{conn}, conf, f)
}

func BindQueues(conn *amqp.Connection, bindings Bindings) error {
	return bindQueues(&binderWrapper{conn}, bindings)
}

type binderWrapper struct {
	*amqp.Connection
}

type connectionWrapper struct {
	*amqp.Connection
}

func (bw *binderWrapper) channel() (binder, error) {
	ch, err := bw.Channel()
	return &channelWrapper{ch}, err
}

func (cw *connectionWrapper) channel() (consumer, error) {
	ch, err := cw.Channel()
	return &channelWrapper{ch}, err
}

func (cw *connectionWrapper) notifyOnClose() chan *amqp.Error {
	closer := make(chan *amqp.Error)
	return cw.NotifyClose(closer)
}

type channelWrapper struct {
	*amqp.Channel
}

func (ch *channelWrapper) close() {
	_ = ch.Close()
}

func (ch *channelWrapper) consume(queue string, autoAck, exclusive, noLocal, noWait bool) (<-chan amqp.Delivery, error) {
	return ch.Consume(queue, "", autoAck, exclusive, noLocal, noWait, nil)
}

func (ch *channelWrapper) declare(d Declare) error {
	_, err := ch.QueueDeclare(d.Name, d.Durable, d.AutoDelete, d.Exclusive, d.NoWait, nil)
	return err
}

func (ch *channelWrapper) bind(b Bind) error {
	return ch.QueueBind(b.Name, b.Topic, b.Exchange, b.NoWait, nil)
}

func (ch *channelWrapper) unbind(u UnBind) error {
	return ch.QueueUnbind(u.Queue, u.Topic, u.Exchange, nil)
}
