package gokafka

// Producer implements kafka producer
type Producer interface {
	Close()
	Produce(topic string, event []byte) (err error)
}

// NewProducer used to new a kafka producer
func NewProducer(broker string) (Producer, error) {
	return newProducer(broker)
}

// Consumer implements kafka consumer
type Consumer interface {
	Close()
	ParseLoop(dealFunc func(b []byte))
	Consume(topic string)
}

// NewConsumer used to new a kafka consumer
func NewConsumer(broker, groupID string) (Consumer, error) {
	return newConsumer(broker, groupID)
}
