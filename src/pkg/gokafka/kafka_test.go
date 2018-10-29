package gokafka

import (
	"fmt"
	"pkg/tcfgs"
	"testing"
	"time"
)

const testTopic = "gotest"

func getKafkaBroker(t *testing.T) string {
	cfg, err := tcfgs.GetTestKafkaCfg()
	if err != nil {
		t.Errorf("get config error : %v", err)
		return "127.0.0.1:9092"
	}
	return fmt.Sprintf("%v:%v", cfg.IP, cfg.Port)
}

func TestProduceToKafka(t *testing.T) {
	p, err := NewProducer(getKafkaBroker(t))
	if err != nil {
		t.Errorf("new a kafka producer failed")
		return
	}
	defer p.Close()

	err = p.Produce(testTopic, []byte("set test"))
	if err != nil {
		t.Errorf("produce event failed")
		return
	}
}

func TestConsumeFromKafka(t *testing.T) {
	s, err := NewConsumer(getKafkaBroker(t), "1")
	if err != nil {
		t.Errorf("new a kafka consumer failed")
		return
	}
	defer s.Close()

	TestProduceToKafka(t)

	dealFunc := func(b []byte) {
		t.Log(string(b))
	}

	go func() {
		go s.ParseLoop(dealFunc)
		s.Consume(testTopic)
		if err != nil {
			t.Errorf("consume event failed")
			return
		}
	}()
	time.Sleep(2 * time.Second)
}
