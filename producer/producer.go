package main

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/Shopify/sarama"
	"tulip.io/hackday/kafka"
)

type MessageEncoder struct {
	WorkerID int    `json:"WorkerID"`
	Type     string `json:"Type"`
	IP       string `json:"IP"`
	Msg      string `json:"Msg"`
}

func (m *MessageEncoder) Encode() ([]byte, error) {
	return json.Marshal(m)
}

func (m *MessageEncoder) Length() int {
	encoded, _ := json.Marshal(m)
	return len(encoded)
}

func producer(cid int) {
	conn, err := sarama.NewSyncProducer([]string{"localhost:9092"}, nil)
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	for {
		msg := kafka.Generate(cid)
		msgEncoder := &MessageEncoder{
			WorkerID: msg.WorkerID,
			Type:     msg.Type,
			IP:       msg.IP,
			Msg:      msg.Msg,
		}
		kmsg := &sarama.ProducerMessage{Topic: "test", Key: sarama.StringEncoder(string(msg.WorkerID)), Value: msgEncoder}

		conn.SendMessage(kmsg)
		time.Sleep(100 * time.Millisecond)
	}
}

func main() {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	args := os.Args[1:]

	if len(args) < 1 {
		log.Panic("Usage: producer <num of producers>")
	}
	numProducers, err := strconv.Atoi(args[0])
	if err != nil {
		log.Panic("Usage: Argument 1 (num of producers) must be an integer")
	}

	for i := 0; i < numProducers; i++ {
		go producer(i)
	}

	go func() {
		<-sigs
		done <- true
	}()

	<-done

}
