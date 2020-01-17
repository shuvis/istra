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

func (cw *connectionWrapper) notifyOnClose() chan error {
	return nil
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

func (ch *channelWrapper) declare(d DeclareConf) error {
	return nil
}

func (ch *channelWrapper) bind(b Bind) error {
	return nil
}

func (ch *channelWrapper) unbind(u UnBind) error {
	return nil
}
