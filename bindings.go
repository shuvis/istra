package istra

type QueueConf struct {
	Name      string
	AutoAck   bool
	Exclusive bool
	NoLocal   bool
	NoWait    bool
}

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

type Bind struct {
	Name     string
	Exchange string
	Queue    string
	Topic    string
	NoWait   bool
}

func (bind Bind) apply(b binder) error {
	return b.bind(bind)
}

type UnBind struct {
	Exchange string
	Queue    string
	Topic    string
}

func (u UnBind) apply(b binder) error {
	return b.unbind(u)
}

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
