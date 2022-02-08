package pkg

import (
	"github.com/Shopify/sarama"
	"time"
)

type KafkaS struct {
	Address []string
}

func (k *KafkaS) PublishMessage(topic string, content []byte) error {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	producer, err := sarama.NewAsyncProducer(k.Address, config)
	if err != nil {
		return err
	}
	defer func(producer sarama.AsyncProducer) {
		producer.AsyncClose()
	}(producer)
	m := &sarama.ProducerMessage{
		Topic:     topic,
		Value:     sarama.ByteEncoder(content),
		Timestamp: time.Time{},
	}
	go func() {
		for _ = range producer.Errors() {
		}
	}()
	producer.Input() <- m
	return nil
}

func (k *KafkaS) Subscribe(topic string, groupId string, f func(*sarama.ConsumerMessage) error) error {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.AutoCommit.Enable = true
	config.Consumer.Offsets.AutoCommit.Interval = 1 * time.Second

	client, err := sarama.NewClient(k.Address, config)
	if err != nil {
		return err
	}
	defer func(client sarama.Client) {
		err := client.Close()
		if err != nil {
		}
	}(client)
	offsetManager, err1 := sarama.NewOffsetManagerFromClient(groupId, client)
	if err1 != nil {
		return err1
	}
	defer func(offsetManager sarama.OffsetManager) {
		err := offsetManager.Close()
		if err != nil {
		}
	}(offsetManager)
	partitionOffsetManager, err2 := offsetManager.ManagePartition(topic, 0)
	if err2 != nil {
		return err2
	}
	defer func(partitionOffsetManager sarama.PartitionOffsetManager) {
		err := partitionOffsetManager.Close()
		if err != nil {
		}
	}(partitionOffsetManager)
	defer offsetManager.Commit()
	consumer, err3 := sarama.NewConsumerFromClient(client)
	if err3 != nil {
		return err3
	}
	nextOffset, _ := partitionOffsetManager.NextOffset()
	pc, _ := consumer.ConsumePartition(topic, 0, nextOffset)
	defer func(pc sarama.PartitionConsumer) {
		err := pc.Close()
		if err != nil {
		}
	}(pc)
	message := pc.Messages()
	for msg := range message {
		err5 := f(msg)
		if err5 != nil {
			partitionOffsetManager.MarkOffset(nextOffset+1, "")
		}
	}
	return nil
}

func NewKafkaS(address []string) *KafkaS {
	k := &KafkaS{Address: address}
	return k
}
