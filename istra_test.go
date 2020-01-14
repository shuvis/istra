package istra

import (
	"errors"
	"reflect"
	"sync"
	"testing"
)

const (
	channel     = "channel"
	channelName = "name"
)

func Test_ProcessQueue(t *testing.T) {

	t.Run("test channel creation call", func(t *testing.T) {
		closeChan := make(chan error)
		amqp := &AmqpMock{Ch: &MessengerChannel{}, CloseChan: closeChan}
		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			ProcessQueue(amqp, QueueConf{}, func(msg []byte) {})
		}()
		closeChan <- errors.New("channel closed")
		wg.Wait()

		want := []string{channel}
		if !reflect.DeepEqual(want, amqp.Calls) {
			t.Errorf("wanted calls %v got %v", want, amqp.Calls)
		}
	})

	t.Run("test channel call return error", func(t *testing.T) {
		amqp := &AmqpChannelErrorMock{}
		assertPanic(t, ErrCreatingChannel, func() {
			ProcessQueue(amqp, QueueConf{}, func(msg []byte) {})
		})
	})

	t.Run("test channel consumed correct config", func(t *testing.T) {
		closeChan := make(chan error)
		channel := ConsumerChannel{}
		amqp := &AmqpMock{Ch: &channel, CloseChan: closeChan}
		conf := QueueConf{Name: channelName, AutoAck: true, Exclusive: true, NoLocal: true, NoWait: true}

		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			ProcessQueue(amqp, conf, func(msg []byte) {})
		}()
		closeChan <- errors.New("channel closed")
		wg.Wait()

		assertConfig(t, channel.Conf, channelName, true, true, true, true)
	})

	t.Run("test channel consumed default config", func(t *testing.T) {
		closeChan := make(chan error)
		channel := ConsumerChannel{}
		amqp := &AmqpMock{Ch: &channel, CloseChan: closeChan}
		conf := QueueConf{Name: channelName}

		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			ProcessQueue(amqp, conf, func(msg []byte) {})
		}()
		closeChan <- errors.New("channel closed")
		wg.Wait()

		assertConfig(t, channel.Conf, channelName, false, false, false, false)
	})

	t.Run("test channel return error", func(t *testing.T) {
		channel := ErrorChannel{}
		amqp := &AmqpMock{Ch: &channel}

		assertPanic(t, ErrConsumingChannel, func() {
			ProcessQueue(amqp, QueueConf{}, func(msg []byte) {})
		})
	})

	t.Run("test channel returning messages", func(t *testing.T) {
		messengerChan := make(chan Messenger)
		closeChan := make(chan error)
		channel := MessengerChannel{MsgChan: messengerChan}
		amqp := &AmqpMock{Ch: &channel, CloseChan: closeChan}
		var result []string
		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			ProcessQueue(amqp, QueueConf{}, func(msg []byte) {
				result = append(result, string(msg))
			})
		}()
		messengerChan <- &MessengerMock{Message: []byte("1")}
		messengerChan <- &MessengerMock{Message: []byte("2")}
		messengerChan <- &MessengerMock{Message: []byte("3")}
		closeChan <- errors.New("channel closed")

		want := []string{"1", "2", "3"}

		wg.Wait()
		if !reflect.DeepEqual(want, result) {
			t.Errorf("wanted calls %v got %v", want, result)
		}
	})
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