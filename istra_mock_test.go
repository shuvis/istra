package istra

type AmqpMock struct {
	Ch        Channel
	Calls     []string
	CloseChan chan error
}

func (a *AmqpMock) Channel() (Channel, error) {
	a.Calls = append(a.Calls, channel)
	return a.Ch, nil
}

func (a *AmqpMock) NotifyOnClose() chan error {
	return a.CloseChan
}

type AmqpChannelErrorMock struct {
}

type MessengerMock struct {
	Message []byte
}

func (m *MessengerMock) Msg() []byte {
	return m.Message
}

type ErrorChannel struct {
}

func (ech *ErrorChannel) Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args interface{}) (<-chan Messenger, error) {
	return nil, ErrConsumingChannel
}

type MessengerChannel struct {
	MsgChan <-chan Messenger
}

func (ch *MessengerChannel) Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args interface{}) (<-chan Messenger, error) {
	return ch.MsgChan, nil
}

type ConsumerChannel struct {
	Conf QueueConf
}

func (ch *ConsumerChannel) Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args interface{}) (<-chan Messenger, error) {
	ch.Conf.Name = queue
	ch.Conf.AutoAck = autoAck
	ch.Conf.Exclusive = exclusive
	ch.Conf.NoLocal = noLocal
	ch.Conf.NoWait = noWait
	return nil, nil
}

func (err *AmqpChannelErrorMock) Channel() (Channel, error) {
	return nil, ErrCreatingChannel
}

func (err *AmqpChannelErrorMock) NotifyOnClose() chan error {
	return nil
}
