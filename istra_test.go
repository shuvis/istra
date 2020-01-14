package istra

import (
	"errors"
	"reflect"
	"sync"
	"testing"
)

const channel = "channel"

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
		conf := QueueConf{Name: "name", AutoAck: true, Exclusive: true, NoLocal: true, NoWait: true}

		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			ProcessQueue(amqp, conf, func(msg []byte) {})
		}()
		closeChan <- errors.New("channel closed")
		wg.Wait()

		if !reflect.DeepEqual(conf, channel.Conf) {
			t.Errorf("wanted calls %v got %v", conf, channel.Conf)
		}
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

func assertPanic(t *testing.T, want error, f func()) {
	defer func() {
		r := recover()
		if r != want {
			t.Errorf("got %q want %q", r, want)
		}
	}()
	f()
}
