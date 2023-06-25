package pubsubs

import (
	"context"
	"log"

	"cloud.google.com/go/pubsub"
	global_config "github.com/vldcreation/sample-cron-go/internal/config"
	"google.golang.org/api/option"
)

const ()

type Pubsubs struct {
	client *pubsub.Client
	config *global_config.Config
}

// NewPubSubs create new instance of pubsub subscriber
func NewPubSubs(ctx context.Context, cfg *global_config.Config) *Pubsubs {
	client, err := pubsub.NewClient(ctx, cfg.PubSub.ProjectID, option.WithCredentialsFile(cfg.PubSub.AccountPath))
	if err != nil {
		log.Fatalf("Could not create pubsub Client: %v", err)
	}

	return &Pubsubs{
		client: client,
		config: cfg,
	}
}
