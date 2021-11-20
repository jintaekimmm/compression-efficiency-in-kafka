package kafka

type Config struct {
	Brokers []string `env:"BROKER" envSeparator:","`
}
