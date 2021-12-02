package bench

import (
	"encoding/json"
	"fmt"
	"github.com/99-66/compression-efficiency-in-kafka/kafka"
	"github.com/Shopify/sarama"
	"log"
	"os"
	"testing"
)

func BenchmarkProducerJsonMarshalingAndGzip(t *testing.B) {
	topicBase := "json-marshaling-and-gzip"

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
			p, err := kafka.NewProducer("gzip")
			if err != nil {
				panic(err)
			}
			defer p.Close()

			for r := range generateJsonSample(file) {
				vJson, err := json.Marshal(r)
				if err != nil {
					log.Printf("item failed marshaling. %v\n", err)
					continue
				}
				msg := &sarama.ProducerMessage{
					Topic: topic,
					Value: sarama.ByteEncoder(vJson),
				}
				p.Input() <- msg
			}
		}()
	}

}

// go test -bench=^BenchmarkProducerJsonMarshalingAndGzip$ . -run=^$ -v -timeout 60m
//goos: darwin
//goarch: amd64
//pkg: github.com/99-66/compression-efficiency-in-kafka/bench
//cpu: Intel(R) Core(TM) i7-8569U CPU @ 2.80GHz
//BenchmarkProducerJsonMarshalingAndGzip
//BenchmarkProducerJsonMarshalingAndGzip-8               1        586673952230 ns/op
//PASS
//ok      github.com/99-66/compression-efficiency-in-kafka/bench    587.095s
