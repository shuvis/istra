package istra

type amqpMock struct {
	ch        channel
	calls     []string
	closeChan chan error
}

func (a *amqpMock) channel() (channel, error) {
	a.calls = append(a.calls, channelMethod)
	return a.ch, nil
}

func (a *amqpMock) notifyOnClose() chan error {
	return a.closeChan
}

type amqpChannelErrorMock struct {
}

type messengerMock struct {
	Message []byte
}

func (m *messengerMock) msg() []byte {
	return m.Message
}

type errorChannel struct {
}

func (ech *errorChannel) consume(queue string, autoAck, exclusive, noLocal, noWait bool) (<-chan messenger, error) {
	return nil, ErrConsumingChannel
}

type messengerChannel struct {
	msgChan <-chan messenger
}

func (ch *messengerChannel) consume(queue string, autoAck, exclusive, noLocal, noWait bool) (<-chan messenger, error) {
	return ch.msgChan, nil
}

type consumerChannel struct {
	Conf QueueConf
}

func (ch *consumerChannel) consume(queue string, autoAck, exclusive, noLocal, noWait bool) (<-chan messenger, error) {
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
