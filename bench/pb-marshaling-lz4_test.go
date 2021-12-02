package bench

import (
	"fmt"
	"github.com/99-66/compression-efficiency-in-kafka/kafka"
	"github.com/Shopify/sarama"
	"log"
	"os"
	"testing"
)

func BenchmarkProducerProtoBufferAndLz4(t *testing.B) {
	topicBase := "pb-marshaling-and-Lz4"

	for _, sample := range samples {
		func() {
			filename := fmt.Sprintf("./samples/%s.csv", sample)
			topic := fmt.Sprintf("%s-%s", topicBase, sample)

			file, err := os.Open(filename)
			if err != nil {
				log.Fatalf("Couldn't open the csv file %s\n", err)
			}
			defer file.Close()

			// Initializing Kafka Producer
			p, err := kafka.NewProducer("lz4")
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
		}()
	}
}
