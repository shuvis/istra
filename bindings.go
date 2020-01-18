package istra

type QueueConf struct {
	Name      string
	AutoAck   bool
	Exclusive bool
	NoLocal   bool
	NoWait    bool
}

type Declare struct {
	Name       string
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
}

func (d Declare) apply(b binder) error {
	return b.declare(d)
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

type DeclareBind struct {
	Declare Declare
	Bind    Bind
}

func (db DeclareBind) apply(b binder) error {
	err := b.declare(db.Declare)
	if err != nil {
		return err
	}
	return b.bind(db.Bind)
}
