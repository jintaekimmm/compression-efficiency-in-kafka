package bench

import (
	"github.com/99-66/compression-efficiency-in-kafka/kafka"
	"github.com/Shopify/sarama"
	"log"
	"os"
	"testing"
)

func BenchmarkProducerProtoBufferAndSnappy(t *testing.B) {
	topic := "pb-marshaling-and-snappy"

	filename := "./samples/data_50m.csv"

	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Couldn't open the csv file %s\n", err)
	}
	defer file.Close()

	// Initializing Kafka Producer
	p, err := kafka.NewProducer("snappy")
	if err != nil {
		panic(err)
	}
	defer p.Close()

	for r := range generateProtoBufferSample(file) {
		if r.Err != nil {
			log.Printf("item failed marshaling. %v\n", err)
			continue
		}
		msg := &sarama.ProducerMessage{
			Topic: topic,
			Value: sarama.ByteEncoder(r.Bytes),
		}
		p.Input() <- msg
	}

}
