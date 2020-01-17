package istra

import (
	"errors"
	"github.com/streadway/amqp"
	"reflect"
	"sync"
	"testing"
)

const (
	channelMethod = "channel"
	channelName   = "name"
)

func Test_ProcessQueue(t *testing.T) {

	t.Run("test channel() is called", func(t *testing.T) {
		closeChan := make(chan error)
		conn := &connectionMock{ch: &messengerChannel{}, closeChan: closeChan}

		wg := processWithWaitingGroup(func() {
			processQueue(conn, QueueConf{}, func(d amqp.Delivery) {})
		})
		closeChan <- errors.New("channel closed")
		wg.Wait()

		want := []string{channelMethod}
		if !reflect.DeepEqual(want, conn.calls) {
			t.Errorf("wanted calls %v got %v", want, conn.calls)
		}
	})

	t.Run("test channel() return error", func(t *testing.T) {
		conn := &amqpChannelErrorMock{}
		assertPanic(t, ErrCreatingChannel, func() {
			processQueue(conn, QueueConf{}, func(d amqp.Delivery) {})
		})
	})

	t.Run("test QueueConf passed to channel.consume(QueueConf)", func(t *testing.T) {
		closeChan := make(chan error)
		channel := consumerChannel{}
		conn := &connectionMock{ch: &channel, closeChan: closeChan}
		conf := QueueConf{Name: channelName, AutoAck: true, Exclusive: true, NoLocal: true, NoWait: true}

		wg := processWithWaitingGroup(func() {
			processQueue(conn, conf, func(d amqp.Delivery) {})
		})
		closeChan <- errors.New("channel closed")
		wg.Wait()

		assertConfig(t, channel.Conf, channelName, true, true, true, true)
	})

	t.Run("test default QueueConf is passed", func(t *testing.T) {
		closeChan := make(chan error)
		channel := consumerChannel{}
		conn := &connectionMock{ch: &channel, closeChan: closeChan}
		conf := QueueConf{Name: channelName}

		wg := processWithWaitingGroup(func() {
			processQueue(conn, conf, func(d amqp.Delivery) {})
		})
		closeChan <- errors.New("channel closed")
		wg.Wait()

		assertConfig(t, channel.Conf, channelName, false, false, false, false)
	})

	t.Run("test consume() return error", func(t *testing.T) {
		channel := errorChannel{}
		conn := &connectionMock{ch: &channel}

		assertPanic(t, ErrConsumingChannel, func() {
			processQueue(conn, QueueConf{}, func(d amqp.Delivery) {})
		})
	})

	t.Run("test deliveries processing", func(t *testing.T) {
		messengerChan := make(chan amqp.Delivery)
		closeChan := make(chan error)
		channel := messengerChannel{msgChan: messengerChan}
		conn := &connectionMock{ch: &channel, closeChan: closeChan}
		var result []string
		wg := processWithWaitingGroup(func() {
			processQueue(conn, QueueConf{}, func(d amqp.Delivery) {
				result = append(result, string(d.Body))
			})
		})
		messengerChan <- amqp.Delivery{Body: []byte("1")}
		messengerChan <- amqp.Delivery{Body: []byte("2")}
		messengerChan <- amqp.Delivery{Body: []byte("3")}
		closeChan <- errors.New("channel closed")

		want := []string{"1", "2", "3"}

		wg.Wait()
		if !reflect.DeepEqual(want, result) {
			t.Errorf("wanted calls %v got %v", want, result)
		}
	})
}

func processWithWaitingGroup(fn func()) *sync.WaitGroup {
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		fn()
	}()
	return &wg
}

func assertConfig(t *testing.T, conf QueueConf, name string, autoAck bool, exclusive bool, noLocal bool, noWait bool) {
	t.Helper()
	if conf.Name != name {
		t.Errorf("wanted Name %v got %v", name, conf.Name)
	}
	if conf.AutoAck != autoAck {
		t.Errorf("wanted AutoAck %v got %v", autoAck, conf.AutoAck)
	}
	if conf.Exclusive != exclusive {
		t.Errorf("wanted Exclusive %v got %v", exclusive, conf.Exclusive)
	}
	if conf.NoLocal != noLocal {
		t.Errorf("wanted NoLocal %v got %v", noLocal, conf.NoLocal)
	}
	if conf.NoWait != noWait {
		t.Errorf("wanted NoWait %v got %v", noWait, conf.NoWait)
	}
}

func assertPanic(t *testing.T, want error, f func()) {
	defer func() {
		r := recover()
		if r != want {
			t.Errorf("got %q want %q", r, want)
		}
	}()
	f()
}
