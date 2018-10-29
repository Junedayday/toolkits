package gokafka

import (
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type producer struct {
	pd *kafka.Producer
}

func newProducer(broker string) (*producer, error) {
	cfgMap := &kafka.ConfigMap{"bootstrap.servers": broker}
	p, err := kafka.NewProducer(cfgMap)
	if err != nil {
		return nil, err
	}
	return &producer{
		pd: p,
	}, nil
}

func (p *producer) Close() {
	p.pd.Close()
}

func (p *producer) Produce(topic string, event []byte) (err error) {
	doneChan := make(chan bool)
	go func() {
		defer close(doneChan)
		for e := range p.pd.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				m := ev
				if m.TopicPartition.Error != nil {
					err = fmt.Errorf(errProduceMessage, m.TopicPartition.Error)
					return
				}
				fmt.Printf("Delivered message to topic %s [%d] at offset %v\n",
					*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
				return

			default:
				fmt.Printf("Ignored event: %s\n", ev)
			}
		}
	}()
	msg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value: []byte(event),
	}
	p.pd.ProduceChannel() <- msg

	_ = <-doneChan
	return nil
}
