package gokafka

import (
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/golang/glog"
)

type kafkaMsg struct {
	data []byte
	err  error
}

type consumer struct {
	cs    *kafka.Consumer
	msgCh chan kafkaMsg
}

func newConsumer(broker, groupID string) (*consumer, error) {
	cfgMap := &kafka.ConfigMap{
		"bootstrap.servers":               broker,
		"group.id":                        groupID,
		"go.events.channel.enable":        true,
		"go.application.rebalance.enable": true,
		"default.topic.config":            kafka.ConfigMap{"auto.offset.reset": "earliest"},
	}
	p, err := kafka.NewConsumer(cfgMap)
	if err != nil {
		return nil, err
	}
	return &consumer{
		cs:    p,
		msgCh: make(chan kafkaMsg),
	}, nil
}

func (c *consumer) Close() {
	c.cs.Close()
}

func (c *consumer) ParseLoop(dealFunc func(b []byte)) {
	for {
		select {
		case msg, ok := <-c.msgCh:
			if ok {
				if msg.err != nil {
					glog.Errorf("parse kafka error : %v", msg.err)
				} else {
					dealFunc(msg.data)
				}
			}
		}
	}
}

func (c *consumer) Consume(topic string) {
	// subscribe topic
	err := c.cs.Subscribe(topic, nil)
	if err != nil {
		c.msgCh <- kafkaMsg{
			err: fmt.Errorf(errSubscribeTopic, topic, err),
		}
		return
	}

	for {
		msg := kafkaMsg{}
		select {
		case ev := <-c.cs.Events():
			switch e := ev.(type) {
			case kafka.AssignedPartitions:
				// fmt.Fprintf(os.Stderr, "%% %v\n", e)
				c.cs.Assign(e.Partitions)
			case kafka.RevokedPartitions:
				// fmt.Fprintf(os.Stderr, "%% %v\n", e)
				c.cs.Unassign()
			case *kafka.Message:
				fmt.Printf("%% Message on %s:\n", e.TopicPartition)
				msg.data = e.Value
				c.msgCh <- msg
			case kafka.PartitionEOF:
				// fmt.Printf("%% Reached %v\n", e)
			case kafka.Error:
				// fmt.Fprintf(os.Stderr, "%% Error: %v\n", e)
				msg.err = fmt.Errorf("kafka error %v", e)
				c.msgCh <- msg
				return
			}
		}
	}
}
