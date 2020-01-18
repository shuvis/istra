package istra

type QueueConf struct {
	Name      string
	AutoAck   bool
	Exclusive bool
	NoLocal   bool
	NoWait    bool
}

type Bindings []interface{}

type Declare struct {
	Name       string
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
}

type Bind struct {
	Exchange string
	Queue    string
	Topic    string
	NoWait   bool
}

type UnBind struct {
	Exchange string
	Queue    string
	Topic    string
}

type DeclareBind struct {
	Declare Declare
	Bind    Bind
}
