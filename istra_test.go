package istra

import (
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
	"reflect"
	"sync"
	"testing"
)

const (
	channelName = "name"
	channel     = "channel()"
	closeMethod = "close()"
	declare     = "testQueue()"
	bind        = "bind()"
	unbind      = "unbind()"
	exchange    = "testExchange()"
)

func Test_consumeQueue(t *testing.T) {

	t.Run("test consumer() is called", func(t *testing.T) {
		closeChan := make(chan *amqp.Error)
		chMock := consumeChannelerMock{ch: &consumerMock{}}
		conn := &connectionMock{consumeChanneler: &chMock, closeChan: closeChan}

		wg := processWithWaitingGroup(func() {
			consumeQueue(conn, QueueConf{}, func(d amqp.Delivery) {})
		})
		closeChan <- &amqp.Error{}
		wg.Wait()

		want := []string{channel}
		if !reflect.DeepEqual(want, chMock.calls) {
			t.Errorf("wanted calls %v got %v", want, chMock.calls)
		}
	})

	t.Run("test consumer() return error", func(t *testing.T) {
		conn := &connectionMock{consumeChanneler: &consumeChannelerMock{err: errors.New("error")}}
		assertPanic(t, "error creating channel: error", func() {
			consumeQueue(conn, QueueConf{}, func(d amqp.Delivery) {})
		})
	})

	t.Run("test QueueConf passed to consumer.consume(QueueConf)", func(t *testing.T) {
		closeChan := make(chan *amqp.Error)
		channel := consumerMock{}
		conn := &connectionMock{consumeChanneler: &consumeChannelerMock{ch: &channel}, closeChan: closeChan}
		conf := QueueConf{Name: channelName, AutoAck: true, Exclusive: true, NoLocal: true, NoWait: true}

		wg := processWithWaitingGroup(func() {
			consumeQueue(conn, conf, func(d amqp.Delivery) {})
		})
		closeChan <- &amqp.Error{}
		wg.Wait()

		assertConfig(t, channel.Conf, channelName, true, true, true, true)
	})

	t.Run("test default QueueConf is passed", func(t *testing.T) {
		closeChan := make(chan *amqp.Error)
		channel := consumerMock{}
		conn := &connectionMock{consumeChanneler: &consumeChannelerMock{ch: &channel}, closeChan: closeChan}
		conf := QueueConf{Name: channelName}

		wg := processWithWaitingGroup(func() {
			consumeQueue(conn, conf, func(d amqp.Delivery) {})
		})
		closeChan <- &amqp.Error{}
		wg.Wait()

		assertConfig(t, channel.Conf, channelName, false, false, false, false)
	})

	t.Run("test consume() return error", func(t *testing.T) {
		conn := &connectionMock{consumeChanneler: &consumeChannelerMock{ch: &consumerMock{err: errors.New("error")}}}

		assertPanic(t, "error consuming testQueue: error", func() {
			consumeQueue(conn, QueueConf{}, func(d amqp.Delivery) {})
		})
	})

	t.Run("test deliveries processing", func(t *testing.T) {
		deliveries := make(chan amqp.Delivery)
		closeChan := make(chan *amqp.Error)
		channel := consumerMock{msgChan: deliveries}
		conn := &connectionMock{consumeChanneler: &consumeChannelerMock{ch: &channel}, closeChan: closeChan}
		var result []string
		wg := processWithWaitingGroup(func() {
			consumeQueue(conn, QueueConf{}, func(d amqp.Delivery) {
				result = append(result, string(d.Body))
			})
		})
		deliveries <- amqp.Delivery{Body: []byte("1")}
		deliveries <- amqp.Delivery{Body: []byte("2")}
		deliveries <- amqp.Delivery{Body: []byte("3")}
		closeChan <- &amqp.Error{}

		want := []string{"1", "2", "3"}

		wg.Wait()
		if !reflect.DeepEqual(want, result) {
			t.Errorf("wanted calls %v got %v", want, result)
		}
	})
}

func Test_QueueBind(t *testing.T) {

	t.Run("test processOperations() returns error", func(t *testing.T) {
		closer := &binderMock{}
		err := processOperations(&bindChannelMock{err: errors.New("error"), b: &binderMock{}}, Actions{})

		wantErr := "error creating channel: error"
		if err == nil || err.Error() != wantErr {
			t.Errorf("wanted '%v' got '%v'", wantErr, err)
		}

		if len(closer.calls) > 0 {
			t.Errorf("didn't expected calls got %v", closer.calls)
		}
	})

	t.Run("test consumer() and close() are called", func(t *testing.T) {
		closer := &binderMock{}
		chMock := &bindChannelMock{b: closer}
		err := processOperations(chMock, Actions{})

		if err != nil {
			t.Errorf("didn't expect error, got '%v'", err)
		}

		wantChannel := []string{channel}
		if !reflect.DeepEqual(wantChannel, chMock.calls) {
			t.Errorf("wanted calls %v got %v", wantChannel, chMock.calls)
		}

		wantCloser := []string{closeMethod}
		if !reflect.DeepEqual(wantCloser, closer.calls) {
			t.Errorf("wanted calls %v got %v", wantCloser, closer.calls)
		}
	})

	t.Run("test processOperations()", func(t *testing.T) {
		d := QueueDeclare{}
		bindErrorQueue := Bind{Exchange: "testExchange", Queue: "errorQueue", Topic: "*.ERROR"}
		bindWarningQueue := Bind{Exchange: "testExchange", Queue: "warningQueue", Topic: "*.WARNING"}
		unBindQueue := UnBind{Exchange: "testExchange", Queue: "infoQueue", Topic: "*.INFO"}
		bindings := Actions{
			d,
			bindErrorQueue,
			bindWarningQueue,
			unBindQueue,
		}

		binder := &binderMock{}
		err := processOperations(&bindChannelMock{b: binder}, bindings)

		if err != nil {
			t.Errorf("didn't expect error, got '%v'", err)
		}

		want := []string{declare, bind, bind, unbind, closeMethod}
		if !reflect.DeepEqual(want, binder.calls) {
			t.Errorf("wanted calls %v got %v", want, binder.calls)
		}
		wantBindings := make([]interface{}, 0)
		wantBindings = append(wantBindings, d, bindErrorQueue, bindWarningQueue, unBindQueue)

		if !reflect.DeepEqual(wantBindings, binder.passedStructs) {
			t.Errorf("wanted passedStructs %v\ngot %v", bindings, binder.passedStructs)
		}
	})

	t.Run("test processOperations() return error", func(t *testing.T) {
		err := errors.New("error")
		dQueue := "d testQueue"
		tests := []struct {
			name       string
			bindings   Actions
			bindErr    error
			declareErr error
			unbindErr  error
			wantCalls  []string
		}{
			{"test testQueue() return error", Actions{QueueDeclare{Name: dQueue}, Bind{}}, nil, err, nil, []string{declare, closeMethod}},
			{"test bind() return error", Actions{QueueDeclare{}, Bind{Queue: "b testQueue"}, Bind{}}, err, nil, nil, []string{declare, bind, closeMethod}},
			{"test unbind() return error", Actions{QueueDeclare{}, Bind{}, UnBind{Queue: "u testQueue"}, Bind{}}, nil, nil, err, []string{declare, bind, unbind, closeMethod}},
		}

		for _, tt := range tests {
			binder := &binderMock{bindErr: tt.bindErr, declareErr: tt.declareErr, unbindErr: tt.unbindErr}
			got := processOperations(&bindChannelMock{b: binder}, tt.bindings)

			if errors.Cause(got) != err {
				t.Errorf("expected error '%v, got '%v'", err, got)
			}

			wantMsg := "binding failed: error"
			if got.Error() != wantMsg {
				t.Errorf("expected error message '%v', got '%v'", wantMsg, got.Error())
			}

			if !reflect.DeepEqual(tt.wantCalls, binder.calls) {
				t.Errorf("wanted calls %v got %v", tt.wantCalls, binder.calls)
			}
		}
	})
}

func processWithWaitingGroup(f func()) *sync.WaitGroup {
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		f()
	}()
	return &wg
}

func assertConfig(t *testing.T, conf QueueConf, name string, autoAck bool, exclusive bool, noLocal bool, noWait bool) {
	t.Helper()
	if conf.Name != name {
		t.Errorf("wanted name %v got %v", name, conf.Name)
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

func assertPanic(t *testing.T, errorMsg string, f func()) {
	defer func() {
		if r, ok := recover().(error); ok {
			if r.Error() != errorMsg {
				t.Errorf("got %q want %q", r, errorMsg)
			}
		}
	}()
	f()
}
