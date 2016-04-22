package main

import (
	"log"
	"os"
	"os/signal"
	"strconv"

	"github.com/Shopify/sarama"
)

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		log.Fatalln("Usage: filteredconsumer <workerid>")
	}

	consumer, err := sarama.NewConsumer([]string{"localhost:9092"}, nil)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := consumer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	partitionConsumer, err := consumer.ConsumePartition("test", 0, sarama.OffsetNewest)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := partitionConsumer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	// Trap SIGINT to trigger a shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	consumed := 0
ConsumerLoop:
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			filter, _ := strconv.Atoi(args[0])
			if msg.Key[0] == byte(filter) {
				log.Printf("Consumed message offset %d\n", msg.Offset)
				log.Println(string(msg.Value))
				consumed++
			}
		case <-signals:
			break ConsumerLoop
		}
	}

	log.Printf("Consumed: %d\n", consumed)
}
