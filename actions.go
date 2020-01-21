package istra

// QueueDeclare configuration for consuming queue with ConsumeQueue() method
type QueueConf struct {
	Name          string
	AutoAck       bool
	Exclusive     bool
	NoLocal       bool
	NoWait        bool
	PrefetchSize  int
	PrefetchCount int
	Global        bool
}

// QueueDeclare action which creates new queue
type QueueDeclare struct {
	Name       string
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
}

func (d QueueDeclare) apply(b binder) error {
	return b.queue(d)
}

// Bind action which binds Queue to Exchange
type Bind struct {
	Queue    string
	Exchange string
	Topic    string
	NoWait   bool
}

func (bind Bind) apply(b binder) error {
	return b.bind(bind)
}

// UnBind action which unbinds Queue from Exchange
type UnBind struct {
	Exchange string
	Queue    string
	Topic    string
}

func (u UnBind) apply(b binder) error {
	return b.unbind(u)
}

// ExchangeDeclare action which creates new exchange
type ExchangeDeclare struct {
	Exchange   string
	Kind       string
	Durable    bool
	AutoDelete bool
	Internal   bool
	NoWait     bool
}

func (ed ExchangeDeclare) apply(b binder) error {
	return b.exchange(ed)
}
