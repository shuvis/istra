package istra

import "amqp"

func New(conn *amqp.Connection) QueueProcessor {
	return &queueProcessor{&connectionWrapper{conn}}
}

type QueueProcessor interface {
	Process(conf QueueConf, f func([]byte))
}

type queueProcessor struct {
	conn *connectionWrapper
}

func (qp *queueProcessor) Process(conf QueueConf, f func([]byte)) {
	processQueue(qp.conn, conf, f)
}

type connectionWrapper struct {
	*amqp.Connection
}

func (cw *connectionWrapper) channel() (channel, error) {
	ch, err := cw.Channel()
	return &channelWrapper{ch}, err
}

func (cw *connectionWrapper) notifyOnClose() chan error {
	return nil
}

type channelWrapper struct {
	*amqp.Channel
}

func (ch *channelWrapper) consume(queue string, autoAck, exclusive, noLocal, noWait bool) (<-chan messenger, error) {
	msg, err := ch.Consume(queue, "", autoAck, exclusive, noLocal, noWait, nil)
	if err != nil {
		return nil, err
	}
	deliveries := make(chan messenger)
	go func() {
		defer close(deliveries)
		for m := range msg {
			deliveries <- &deliveryWrapper{m}
		}
	}()
	return deliveries, nil
}

type deliveryWrapper struct {
	amqp.Delivery
}

func (d *deliveryWrapper) msg() []byte {
	return d.Body
}
