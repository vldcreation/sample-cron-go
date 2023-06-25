package pubsubs

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/pubsublite/pscompat"
	"github.com/rs/zerolog/log"
)

type gPublisher struct {
	client *pubsub.Client
}

// NewGPublisher create new instance of pubsub publisher
func NewGPublisher(p *Pubsubs) *gPublisher {
	return &gPublisher{client: p.client}
}

// Publish publish message to the topic
func (p *gPublisher) Publish(ctx context.Context, msg *Message) error {
	tp := createTopicIfNotExists(p.client, msg.Topic)

	payload := &pubsub.Message{
		Data:        msg.Data,
		Attributes:  msg.Attribute,
		PublishTime: time.Now(),
	}
	result := tp.Publish(ctx, payload)

	//_ = result
	_, err := result.Get(ctx)
	return err
}

func createTopicIfNotExists(c *pubsub.Client, topic string) *pubsub.Topic {
	ctx := context.Background()
	t := c.Topic(topic)
	ok, err := t.Exists(ctx)
	// if err != nil {
	// 	log.Fatal().Err(err)
	// }
	if ok {
		return t
	}
	t, err = c.CreateTopic(ctx, topic)
	if err != nil {
		log.Fatal().Msgf("Failed to create the topic: %v", err)
	}
	return t
}

type gSubscriber struct {
	client *pubsub.Client
}

// NewGSubscriber create new instance of pubsub subscriber
func NewGSubscriber(p *Pubsubs) *gSubscriber {
	return &gSubscriber{client: p.client}
}

// Subscribe publish message to the topic
func (p *gSubscriber) Subscribe(ctx context.Context, handler func(context.Context, *Message)) error {
	defer p.client.Close()

	cfg := defaults()

	log.Printf("Subscribing to topic %s", cfg.Topic)

	t := createTopicIfNotExists(p.client, cfg.Topic)

	sub, err := createSubsIfNotExists(p.client, cfg.Topic, t)
	if err != nil {
		log.Fatal().Msgf("Failed to create the topic: %v", err)
		return err
	}

	sub.ReceiveSettings.Synchronous = cfg.SubscribeAsync
	sub.ReceiveSettings.MaxOutstandingMessages = cfg.MaxConcurrent

	var (
		// mu protects the received counter.
		// NOTE: This is not strictly necessary if you are only accessing the counter
		received int32
	)

	err = sub.Receive(ctx, func(xCtx context.Context, msg *pubsub.Message) {
		handler(xCtx, &Message{
			ID:   msg.ID,
			Data: msg.Data,
		})
		// NOTE: May be called concurrently; synchronize access to shared memory.
		atomic.AddInt32(&received, 1)

		// Metadata decoded from the message ID contains the partition and offset.
		metadata, err := pscompat.ParseMessageMetadata(msg.ID)
		if err != nil {
			log.Fatal().Msgf("Failed to parse %q: %v", msg.ID, err)
		}

		fmt.Printf("Received (partition=%d, offset=%d): %s\n", metadata.Partition, metadata.Offset, string(msg.Data))

		msg.Ack()
	})

	if err != nil {
		return fmt.Errorf("sub.Receive: %v", err)
	}

	log.Info().Msg(fmt.Sprintf("%s received %d messages\n", cfg.Topic, received))
	return nil

}

func createSubsIfNotExists(client *pubsub.Client, name string, topic *pubsub.Topic) (*pubsub.Subscription, error) {
	ctx := context.Background()
	sub := client.Subscription(name)
	ok, err := sub.Exists(ctx)
	if err != nil {
		return nil, err
	}

	if ok {
		return sub, nil
	}

	// [START pubsub_create_pull_subscription]
	sub, err = client.CreateSubscription(ctx, name, pubsub.SubscriptionConfig{
		Topic:       topic,
		AckDeadline: 20 * time.Second,
	})
	if err != nil {
		return nil, err
	}
	fmt.Printf("Created subscription: %v\n", sub)
	// [END pubsub_create_pull_subscription]
	return sub, nil
}
