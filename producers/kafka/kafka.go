package kafka

import "github.com/Shopify/sarama"

// Producer for kafka messages
type Producer struct {
	pr sarama.AsyncProducer
}

// NewProducer for Kafka
func NewProducer(connection string) (Producer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	producer, err := sarama.NewAsyncProducer([]string{connection}, config)
	if err != nil {
		return Producer{}, err
	}
	p := Producer{
		pr: producer,
	}
	return p, nil
}

func (p Producer) publish(msg string, topic string) {
	message := &sarama.ProducerMessage{Topic: topic, Value: sarama.StringEncoder(msg)}
	p.pr.Input() <- message
}
