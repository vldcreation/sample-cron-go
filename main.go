package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/vldcreation/sample-cron-go/internal/app"
	"github.com/vldcreation/sample-cron-go/internal/pubsubs"
	"github.com/vldcreation/sample-cron-go/internal/storage"
	"github.com/vldcreation/sample-cron-go/internal/utils"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	initApp := app.App{}

	app.Run(ctx, &initApp)

	// init instance
	// hashs := hashs.NewHashs(sha256.New224(), false, nil)

	// declare path file
	pathFile := "./test_data/sample.jpeg"

	// check file exist or not
	ok, err := utils.FileExists(pathFile)
	if err != nil {
		log.Fatalf("error check file exist: %v\n", err)
		panic(err)
	}

	if !ok {
		log.Fatalf("file not exist: %v\n", pathFile)
		panic(err)
	}

	fileName, ext, err := utils.ParseFile(pathFile)
	if err != nil {
		log.Fatalf("error parse file: %v\n", err)
		panic(err)
	}

	log.Printf("name: %v, ext: %v\n", fileName, ext)

	// read path file

	// put file to storage
	err = initApp.Storage.FPut(ctx, initApp.Config.Storage.Bucket, fileName+"."+ext, pathFile, true, "image/jpeg")
	if err != nil {
		log.Fatalf("error put file to storage: %v\n", err)
		panic(err)
	}

	var URLGenerated string

	// first generated url
	bt, err := initApp.Storage.Get(ctx, initApp.Config.Storage.Bucket, fileName+"."+ext)
	if err != nil {
		log.Fatalf("error get file from storage: %v\n", err)
		panic(err)
	}

	URLGenerated = string(bt)
	log.Printf("first url: %v\n", URLGenerated)

	//
	// initialize adjustment to stop cron
	//
	var (
		maxCron int = 2
		curCron int = 0
		// use as flat channel to stop cron
		isStop = make(chan bool)
	)

	fnIter := func() {
		if curCron >= maxCron {
			isStop <- true
			// instead of using channel, we can publish message to pubsub (just for sample purpose)
			// to stop cron
			log.Printf("Send message to stop cron\n")
			if err := initApp.Publisherer.Publish(ctx, &pubsubs.Message{
				Topic: initApp.Config.PubSub.Topic,
				Data:  []byte("stop cron"),
			}); err != nil {
				log.Fatalf("error publish message: %v\n", err)
			}
		} else {
			log.Printf("current cron: %v\n", curCron)
			log.Printf("remain cron: %d\n", maxCron-curCron)
			curCron++
		}
	}

	// declare function to stop cron
	fnStop := func(ctx context.Context, message *pubsubs.Message) {
		log.Print(fmt.Sprintf("%s received %s\n", message.Topic, string(message.Data)))
		initApp.Cron.Stop()
		os.Exit(0)
	}

	// declare function to reSigned url, then send to cron
	ResignedURLFunc := func() {
		go fnIter()
		log.Println("reSigned url for file ", fileName)
		resResignedUrl, err := initApp.Storage.ReSignedURL(ctx, initApp.Config.Storage.Bucket, fileName+"."+ext, URLGenerated)
		if err != nil {
			log.Fatalf("error reSigned url: %v\n", err)
			panic(err)
		}

		URLGenerated = resResignedUrl
		log.Printf("new reSigned url: %v\n", resResignedUrl)

		//
		// initialize subcriber to stop cron
		//
		go func() {
			select {
			case <-isStop:
				if err := initApp.Subscriberer.Subscribe(ctx, fnStop); err != nil {
					log.Fatalf("error subscribe: %v\n", err)
				}
			}
		}()
	}

	if err := initApp.Cron.AddJobWithInterval(storage.Test10Seconds, ResignedURLFunc); err != nil {
		log.Fatalf("error add job: %v\n", err)
		panic(err)
	}

	time.Sleep(storage.Test10Seconds)
	initApp.Cron.StartWithBlocking()

	// defer initApp.Cron.Stop()
}
