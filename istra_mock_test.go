package istra

import "github.com/streadway/amqp"

type connectionMock struct {
	ch        channel
	calls     []string
	closeChan chan error
}

func (a *connectionMock) channel() (channel, error) {
	a.calls = append(a.calls, channelMethod)
	return a.ch, nil
}

func (a *connectionMock) notifyOnClose() chan error {
	return a.closeChan
}

type amqpChannelErrorMock struct {
}

type errorChannel struct {
}

func (ech *errorChannel) consume(queue string, autoAck, exclusive, noLocal, noWait bool) (<-chan amqp.Delivery, error) {
	return nil, ErrConsumingChannel
}

type deliveryChannel struct {
	msgChan <-chan amqp.Delivery
}

func (ch *deliveryChannel) consume(queue string, autoAck, exclusive, noLocal, noWait bool) (<-chan amqp.Delivery, error) {
	return ch.msgChan, nil
}

type consumerChannel struct {
	Conf QueueConf
}

func (ch *consumerChannel) consume(queue string, autoAck, exclusive, noLocal, noWait bool) (<-chan amqp.Delivery, error) {
	ch.Conf.Name = queue
	ch.Conf.AutoAck = autoAck
	ch.Conf.Exclusive = exclusive
	ch.Conf.NoLocal = noLocal
	ch.Conf.NoWait = noWait
	return nil, nil
}

func (err *amqpChannelErrorMock) channel() (channel, error) {
	return nil, ErrCreatingChannel
}

func (err *amqpChannelErrorMock) notifyOnClose() chan error {
	return nil
}
