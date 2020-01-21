[![Build Status](https://travis-ci.org/shuvis/istra.svg)](http://travis-ci.org/shuvis/istra) 
[![Coverage Status](https://coveralls.io/repos/github/shuvis/istra/badge.svg?branch=master)](https://coveralls.io/github/shuvis/istra?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/shuvis/istra)](https://goreportcard.com/report/github.com/shuvis/istra)
[![GoDoc](https://godoc.org/github.com/shuvis/istra?status.svg)](http://godoc.org/github.com/shuvis/istra)
[![MIT licensed](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/shuvis/istra/master/LICENSE)

# Go AMQP wrapper

This is an [AMQP Library](https://github.com/streadway/amqp) wrapper.

## Example

```go

import (
	"github.com/shuvis/istra"
	"github.com/streadway/amqp"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}

	err = istra.Process(conn,
		istra.ExchangeDeclare{Exchange: "testExchange", Kind: "fanout", Durable: true, NoWait: false},
		istra.QueueDeclare{Name: "testQueue", Durable: true, AutoDelete: true, Exclusive: true},
		istra.Bind{Queue: "testQueue", Exchange: "testExchange", Topic: "topic", NoWait: true})

	if err != nil {
		panic(err)
	}

	go istra.ConsumeQueue(conn, istra.QueueConf{Name: "testQueue"}, func(d amqp.Delivery) {
		// process delivery
	})

	err = istra.Process(conn,
		istra.UnBind{Queue: "testQueue", Exchange: "testExchange", Topic: "topic"})

	if err != nil {
		panic(err)
	}
}

```

## License

MIT License - see LICENSE for more details.

