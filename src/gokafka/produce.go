package gokafka

import (
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/golang/glog"
)

const (
	errProduceMessage = "Delivery failed: %v"
	errIgnoredEvent   = "Produce an ignored event %v"
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
	return &producer{pd: p}, nil
}

func (p *producer) Close() {
	p.pd.Close()
}

func (p *producer) ProduceToTopic(topic string, event []byte) (err error) {
	doneChan := make(chan struct{})
	go func() {
		defer close(doneChan)
		err = p.produceEvent()
	}()

	p.pd.ProduceChannel() <- &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value: []byte(event),
	}

	_ = <-doneChan
	return nil
}

func (p *producer) produceEvent() (err error) {
	for e := range p.pd.Events() {
		switch ev := e.(type) {
		case *kafka.Message:
			if ev.TopicPartition.Error != nil {
				err = fmt.Errorf(errProduceMessage, ev.TopicPartition.Error)
				return
			}
			glog.V(1).Infof("Delivered message to topic %s [%d] at offset %v\n",
				*ev.TopicPartition.Topic, ev.TopicPartition.Partition, ev.TopicPartition.Offset)
		default:
			err = fmt.Errorf(errIgnoredEvent, e)
		}
		return
	}
	return
}
