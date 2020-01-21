package istra

import "github.com/streadway/amqp"

type connectionMock struct {
	consumeChanneler
	closeChan chan *amqp.Error
}

func (cm *connectionMock) notifyOnClose() chan *amqp.Error {
	return cm.closeChan
}

type consumeChannelerMock struct {
	ch    consumer
	err   error
	calls []string
}

func (ccm *consumeChannelerMock) channel() (consumer, error) {
	ccm.calls = append(ccm.calls, channel)
	return ccm.ch, ccm.err
}

type consumerMock struct {
	Conf    QueueConf
	err     error
	msgChan <-chan amqp.Delivery
}

func (cm *consumerMock) consume(queue string, autoAck, exclusive, noLocal, noWait bool) (<-chan amqp.Delivery, error) {
	cm.Conf.Name = queue
	cm.Conf.AutoAck = autoAck
	cm.Conf.Exclusive = exclusive
	cm.Conf.NoLocal = noLocal
	cm.Conf.NoWait = noWait
	return cm.msgChan, cm.err
}

type operatorChannelMock struct {
	b     operator
	err   error
	calls []string
}

func (bcm *operatorChannelMock) channel() (operator, error) {
	bcm.calls = append(bcm.calls, channel)
	return bcm.b, bcm.err
}

type operatorMock struct {
	bindErr       error
	declareErr    error
	unbindErr     error
	calls         []string
	passedStructs []interface{}
}

func (bm *operatorMock) bind(b Bind) error {
	bm.passedStructs = append(bm.passedStructs, b)
	bm.calls = append(bm.calls, bind)
	return bm.bindErr
}

func (bm *operatorMock) queue(conf QueueDeclare) error {
	bm.passedStructs = append(bm.passedStructs, conf)
	bm.calls = append(bm.calls, declare)
	return bm.declareErr
}

func (bm *operatorMock) unbind(u UnBind) error {
	bm.passedStructs = append(bm.passedStructs, u)
	bm.calls = append(bm.calls, unbind)
	return bm.unbindErr
}

func (bm *operatorMock) exchange(ed ExchangeDeclare) error {
	bm.passedStructs = append(bm.passedStructs, ed)
	bm.calls = append(bm.calls, exchange)
	return bm.unbindErr
}

func (bm *operatorMock) close() {
	bm.calls = append(bm.calls, closeMethod)
}
