// +build integration

package istra

import (
	"github.com/streadway/amqp"
	"reflect"
	"testing"
)

var dsn = "amqp://guest:guest@localhost:5672/"

func Test_Integration(t *testing.T) {
	conn, err := amqp.Dial(dsn)
	if err != nil {
		t.Errorf("Could not connect to RabbitMQ: %v", err.Error())
		return
	}
	var messages []string
	ConsumeQueue(conn, QueueConf{Name: "queue"}, func(d amqp.Delivery) {
		messages = append(messages, string(d.Body))
	})

	want := []string{"1", "2", "3"}
	if !reflect.DeepEqual(want, messages) {
		t.Errorf("wanted calls %v got %v", want, messages)
	}
}
