// +build integration

package istra

import (
	"github.com/streadway/amqp"
	"os"
	"reflect"
	"testing"
)

const (
	testQueue    = "testQueue"
	testExchange = "testExchange"
	amqpUrl      = "AMQP_URL"
)

func Test_Integration(t *testing.T) {
	url := os.Getenv(amqpUrl)
	if url == "" {
		t.Fatalf("environment variable 'AMQP_URL' is not set")
	}

	conn, err := amqp.Dial(url)
	if err != nil {
		t.Fatalf("didn't expected error, got '%v'", err)
	}

	err = ProcessOperations(conn, Operations{
		ExchangeDeclare{Exchange: testExchange, Kind: "fanout"},
		QueueDeclare{Name: testQueue},
		Bind{Queue: testQueue, Exchange: testExchange}})
	if err != nil {
		t.Fatalf("didn't expected error, got '%v'", err)
	}

	var messages []string
	go ConsumeQueue(conn, QueueConf{Name: testQueue, AutoAck: true}, func(d amqp.Delivery) {
		messages = append(messages, string(d.Body))
	})

	ch, err := conn.Channel()
	if err != nil {
		t.Fatalf("didn't expected error, got '%v'", err)
	}

	sendToExchange(ch, testExchange, "1")
	sendToExchange(ch, testExchange, "2")
	sendToExchange(ch, testExchange, "3")

	err = ProcessOperations(conn, Operations{UnBind{Queue: testQueue, Exchange: testExchange}})
	if err != nil {
		t.Fatalf("didn't expected error, got '%v'", err)
	}

	sendToExchange(ch, testExchange, "4")

	err = ch.Close()
	if err != nil {
		t.Fatalf("didn't expected error, got '%v'", err)
	}

	err = conn.Close()
	if err != nil {
		t.Fatalf("didn't expected error, got '%v'", err)
	}

	want := []string{"1", "2", "3"}
	if !reflect.DeepEqual(want, messages) {
		t.Errorf("wanted calls %v got %v", want, messages)
	}
}

func sendToExchange(ch *amqp.Channel, e string, msg string) {
	_ = ch.Publish(e, "", false, false, amqp.Publishing{ContentType: "text/plain", Body: []byte(msg)})
}
