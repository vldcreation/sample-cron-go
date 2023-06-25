package app

import (
	"context"
	"log"

	"github.com/vldcreation/sample-cron-go/internal/config"
	cron_jobs "github.com/vldcreation/sample-cron-go/internal/cron-jobs"
	"github.com/vldcreation/sample-cron-go/internal/pubsubs"
	"github.com/vldcreation/sample-cron-go/internal/storage"
)

type App struct {
	Config       *config.Config
	Cron         *cron_jobs.Cron
	Storage      storage.Storage
	Publisherer  pubsubs.Publisher
	Subscriberer pubsubs.Subscriberer
}

func Run(ctx context.Context, app *App) {
	// init config
	conf := config.NewAppConfig()
	app.Config = conf

	log.Printf("app config: %+v\n", app.Config.Storage)

	// init cron
	app.Cron = cron_jobs.NewCron()

	// init pubsub
	initPubsubs := pubsubs.NewPubSubs(ctx, conf)
	app.Publisherer = pubsubs.NewGPublisher(initPubsubs)
	app.Subscriberer = pubsubs.NewGSubscriber(initPubsubs)

	// init storage
	// default storage
	defaultStorage, err := storage.NewMinio()
	if err != nil {
		log.Fatalf("error init minio: %v\n", err)
		panic(err)
	}
	// switcher
	switch conf.App.APP_ENV {
	case "dev":
		app.Storage = defaultStorage
		break
	case "prod":
		app.Storage, err = storage.NewGCS(ctx, conf.GCS.AccountPath)
		if err != nil {
			log.Fatalf("error init gcs: %v\n", err)
			panic(err)
		}
		break
	default:
		app.Storage = defaultStorage
		break
	}

	log.Println("app initialized successfully")
}
