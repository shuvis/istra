package istra

import "github.com/streadway/amqp"

func ProcessQueue(conn *amqp.Connection, conf QueueConf, f func(amqp.Delivery)) {
	processQueue(&connectionWrapper{conn}, conf, f)
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

func (ch *channelWrapper) consume(queue string, autoAck, exclusive, noLocal, noWait bool) (<-chan amqp.Delivery, error) {
	return ch.Consume(queue, "", autoAck, exclusive, noLocal, noWait, nil)
}
