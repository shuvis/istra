package istra

type QueueConf struct {
	Name      string
	AutoAck   bool
	Exclusive bool
	NoLocal   bool
	NoWait    bool
}

type Binding struct {
	Exchange string
	Queue    string
	Topic    string
}
