package istra

import "github.com/streadway/amqp"

// ConsumeQueue calls handler function on each message delivered to a testQueue
func ConsumeQueue(conn *amqp.Connection, conf QueueConf, f func(amqp.Delivery)) {
	consumeQueue(&consumerWrapper{conn, conf}, conf, f)
}

// ProcessOperations process passed bindings
func ProcessOperations(conn *amqp.Connection, operations Operations) error {
	return processOperations(&operatorWrapper{conn}, operations)
}

type operatorWrapper struct {
	*amqp.Connection
}

type consumerWrapper struct {
	*amqp.Connection
	QueueConf
}

func (bw *operatorWrapper) channel() (operator, error) {
	ch, err := bw.Channel()
	return &channelWrapper{ch}, err
}

func (cw *consumerWrapper) channel() (consumer, error) {
	ch, err := cw.Channel()
	if err != nil {
		return nil, err
	}
	err = ch.Qos(cw.PrefetchCount, cw.PrefetchSize, cw.Global)
	return &channelWrapper{ch}, err
}

func (cw *consumerWrapper) notifyOnClose() chan *amqp.Error {
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

func (ch *channelWrapper) queue(d QueueDeclare) error {
	_, err := ch.QueueDeclare(d.Name, d.Durable, d.AutoDelete, d.Exclusive, d.NoWait, nil)
	return err
}

func (ch *channelWrapper) bind(b Bind) error {
	return ch.QueueBind(b.Queue, b.Topic, b.Exchange, b.NoWait, nil)
}

func (ch *channelWrapper) unbind(u UnBind) error {
	return ch.QueueUnbind(u.Queue, u.Topic, u.Exchange, nil)
}

func (ch *channelWrapper) exchange(ed ExchangeDeclare) error {
	return ch.ExchangeDeclare(ed.Exchange, ed.Kind, ed.Durable, ed.AutoDelete, ed.Internal, ed.NoWait, nil)
}
