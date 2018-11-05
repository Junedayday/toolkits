package gokafka

import (
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/golang/glog"
)

const (
	errSubscribeTopic      = "Subscribe topic %v error : %v"
	errConsumeChanClosed   = "Conusmer message channel is closed!"
	errConsumePartitionEOF = "kafka partition EOF error %v"
	errConsumeKafka        = "kafka consume with kafka error %v"
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
	c, err := kafka.NewConsumer(cfgMap)
	if err != nil {
		return nil, err
	}
	return &consumer{
		cs:    c,
		msgCh: make(chan kafkaMsg),
	}, nil
}

func (c *consumer) Close() {
	c.cs.Close()
}

func (c *consumer) DealLoop(dealFunc func(data []byte)) {
	for {
		select {
		case msg, ok := <-c.msgCh:
			if ok {
				if msg.err != nil {
					glog.Errorf("kafka message has error : %v", msg.err)
				} else {
					dealFunc(msg.data)
				}
			} else {
				glog.Warningf(errConsumeChanClosed)
			}
		}
	}
}

func (c *consumer) ConsumeFromTopic(topic string) {
	err := c.cs.Subscribe(topic, nil)
	defer close(c.msgCh)
	if err != nil {
		c.msgCh <- kafkaMsg{err: fmt.Errorf(errSubscribeTopic, topic, err)}
		return
	}

	for {
		msg := kafkaMsg{}
		select {
		case ev := <-c.cs.Events():
			switch e := ev.(type) {
			case kafka.AssignedPartitions:
				c.cs.Assign(e.Partitions)
			case kafka.RevokedPartitions:
				c.cs.Unassign()
			case *kafka.Message:
				glog.Infof("Message consumed on %s:", e.TopicPartition)
				msg.data = e.Value
				c.msgCh <- msg
			case kafka.PartitionEOF:
				msg.err = fmt.Errorf(errConsumePartitionEOF, e)
				c.msgCh <- msg
				return
			case kafka.Error:
				msg.err = fmt.Errorf(errConsumeKafka, e)
				c.msgCh <- msg
				return
			}
		}
	}
}
