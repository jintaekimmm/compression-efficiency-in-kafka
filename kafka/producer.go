package kafka

import (
	"errors"
	"github.com/Shopify/sarama"
	"github.com/caarlos0/env"
	"time"
)

func newAsyncProducer(conf *Config, compressionType string) (sarama.AsyncProducer, error) {
	saramaCfg := sarama.NewConfig()
	saramaCfg.Producer.Flush.Frequency = time.Second * 3

	switch compressionType {
	case "gzip":
		saramaCfg.Producer.Compression = sarama.CompressionGZIP
	case "snappy":
		saramaCfg.Producer.Compression = sarama.CompressionSnappy
	case "lz4":
		saramaCfg.Producer.Compression = sarama.CompressionLZ4
	case "zstd":
		// zstd compression requires Version >= V2_1_0_0
		saramaCfg.Version = sarama.V2_7_0_0
		saramaCfg.Producer.Compression = sarama.CompressionZSTD
	case "":
	}
	producer, err := sarama.NewAsyncProducer(conf.Brokers, saramaCfg)
	if err != nil {
		return nil, err
	}

	return producer, nil
}

func NewProducer(compressionType string) (sarama.AsyncProducer, error) {
	conf := Config{}
	if err := env.Parse(&conf); err != nil {
		return nil, errors.New("cloud not load kafka environment variables")
	}

	return newAsyncProducer(&conf, compressionType)
}
