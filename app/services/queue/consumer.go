package queue

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/wagslane/go-rabbitmq"

	"github.com/yanchevdimitar/RSS-Reader-Service/app/database"
)

type DefaultConsumer struct {
	Queue
	rssRepo database.RSSRepository
}

func NewDefaultConsumer(rssRepo database.RSSRepository) Processor {
	return DefaultConsumer{NewDefaultQueue(), rssRepo}
}

func (c DefaultConsumer) Process() {
	consumer, err := c.Create()
	if err != nil {
		log.Fatal(err)
	}

	defer consumer.Close()

	interval, _ := strconv.Atoi(os.Getenv("QUEUE_CONSUMER_INTERVAL"))
	ticker := time.NewTicker(time.Duration(interval) * time.Second)

	done := c.GracefullyStop()
	for {
		select {
		case <-ticker.C:
			err = consumer.StartConsuming(
				func(d rabbitmq.Delivery) rabbitmq.Action {
					var rss database.RSS
					urls := strings.Split(string(d.Body), ",")
					c.rssRepo.DeleteAll()

					for i := range urls {
						rss = database.RSS{URL: strings.ReplaceAll(urls[i], " ", "")}
						result := c.rssRepo.Create(rss)
						if result.Error != nil {
							log.Printf("Consumed error: %v", result.Error)
						}
					}
					return rabbitmq.Ack
				},
				os.Getenv("QUEUE_CONSUMER"),
				[]string{""},
				rabbitmq.WithConsumeOptionsConcurrency(1),
				rabbitmq.WithConsumeOptionsQueueDurable,
				rabbitmq.WithConsumeOptionsConsumerAutoAck(true),
			)
			log.Println("Consumer...")
			if err != nil {
				log.Println(err)
			}
		case <-done:
			fmt.Println("stopping consumer")
			return
		}
	}
}

func (c DefaultConsumer) Create() (consumer rabbitmq.Consumer, err error) {
	consumer, err = rabbitmq.NewConsumer(
		"amqp://"+os.Getenv("RB_QUEUE_USER")+":"+os.Getenv("RB_QUEUE_PASS")+"@localhost:5672/",
		rabbitmq.Config{},
		rabbitmq.WithConsumerOptionsLogging,
	)

	return
}
