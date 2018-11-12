package gokafka

// Producer implements kafka producer
type Producer interface {
	Close()
	ProduceToTopic(topic string, event []byte) (err error)
}

// NewProducer : broker is the url for kafka like "127.0.0.1:9092"
func NewProducer(broker string) (Producer, error) {
	return newProducer(broker)
}

// Consumer implements kafka consumer
type Consumer interface {
	Close()
	DealLoop(dealFunc func(data []byte))
	ConsumeFromTopic(topic string)
}

// NewConsumer : broker is the url for kafka like "127.0.0.1:9092", groupID is the identify of the consumer
func NewConsumer(broker, groupID string) (Consumer, error) {
	return newConsumer(broker, groupID)
}
