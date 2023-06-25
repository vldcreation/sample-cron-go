package pubsubs

import (
	global_config "github.com/vldcreation/sample-cron-go/internal/config"
)

type Option func(cfg *config)

type config struct {
	MaxConcurrent  int
	SubscribeAsync bool
	Topic          string
}

func defaults() *config {
	cfg := global_config.NewAppConfig().PubSub

	return &config{
		SubscribeAsync: false,
		Topic:          cfg.Topic,
	}
}

func WithTopic(v string) Option {
	return func(c *config) {
		c.Topic = v
	}
}

func WithMaxConcurrent(v int) Option {
	return func(c *config) {
		c.MaxConcurrent = v
	}
}

func WithSubscribeAsync(v bool) Option {
	return func(c *config) {
		c.SubscribeAsync = v
	}
}
