package queue

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/wagslane/go-rabbitmq"
	reader "github.com/yanchevdimitar/go-rss-reader"

	"github.com/yanchevdimitar/RSS-Reader-Service/app/database"
)

type DefaultPublisher struct {
	Queue
	rssRepo database.RSSRepository
}

func NewDefaultPublisher(rssRepo database.RSSRepository) Processor {
	return DefaultPublisher{NewDefaultQueue(), rssRepo}
}

func (c DefaultPublisher) Process() {
	publisher, err := c.Create()

	if err != nil {
		log.Fatal(err)
	}

	defer publisher.Close()

	interval, _ := strconv.Atoi(os.Getenv("QUEUE_PUBLISHER_INTERVAL"))
	ticker := time.NewTicker(time.Duration(interval) * time.Second)

	done := c.GracefullyStop()
	for {
		select {
		case <-ticker.C:
			var urls []string
			result, rss := c.rssRepo.Get()
			if result.Error != nil {
				log.Printf("Consumed error: %v", err)
			}

			for url := range rss {
				urls = append(urls, rss[url].URL)
			}

			rssItems, _ := reader.NewRssReader(urls, reader.NewRssParser()).Parse()
			feed, _ := json.Marshal(rssItems)

			err = publisher.Publish(
				feed,
				[]string{os.Getenv("QUEUE_PUBLISHER")},
				rabbitmq.WithPublishOptionsContentType("application/json"),
			)
			log.Println("Publisher...")
			if err != nil {
				log.Println(err)
			}
		case <-done:
			fmt.Println("stopping publisher")
			return
		}
	}
}

func (c DefaultPublisher) Create() (publisher *rabbitmq.Publisher, err error) {
	publisher, err = rabbitmq.NewPublisher(
		"amqp://"+os.Getenv("RB_QUEUE_USER")+":"+os.Getenv("RB_QUEUE_PASS")+"@"+os.Getenv("RB_QUEUE_HOST"),
		rabbitmq.Config{},
		rabbitmq.WithPublisherOptionsLogging,
	)

	return
}
