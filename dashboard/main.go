package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Shopify/sarama"
	"github.com/jroimartin/gocui"
	"time"
	"tulip.io/hackday/kafka"
)

func main() {
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

	gui := gocui.NewGui()
	if err := gui.Init(); err != nil {
		panic(err)
	}
	defer gui.Close()

	gui.SetLayout(mainScreenTurnOn)
	gui.SetKeybinding("", 'q', gocui.ModNone, quit)

	done := make(chan bool)
	names := make([]string, 10)

	go streamKafka(gui, partitionConsumer, done)
	go updateNames(partitionConsumer, &names, done)
	go displayNames(gui, &names)

	if err := gui.MainLoop(); err != nil && err != gocui.ErrQuit {
		panic(err)
	}
}

func mainScreenTurnOn(gui *gocui.Gui) error {
	x, y := gui.Size()
	v, err := gui.SetView("stream", 0, 0, x, 10)
	if err != nil && err != gocui.ErrUnknownView {
		panic(err)
	}
	v.Autoscroll = true

	v, err = gui.SetView("mainscreen", 10, 10, x-10, y-10)
	if err != nil && err != gocui.ErrUnknownView {
		panic(err)
	}

	v.Clear()
	return nil
}

func quit(gui *gocui.Gui, view *gocui.View) error {
	return gocui.ErrQuit
}

func streamKafka(gui *gocui.Gui, partitionConsumer sarama.PartitionConsumer, done <-chan bool) {
ConsumerLoop:
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			gui.Execute(func(g *gocui.Gui) error {
				v, err := g.View("stream")
				if err != nil {
					panic(err)
				}
				fmt.Fprintln(v, string(msg.Value))
				return nil
			})
		case <-done:
			break ConsumerLoop
		}
	}
}

func updateNames(partitionConsumer sarama.PartitionConsumer, pNames *[]string, done <-chan bool) {
ConsumerLoop:
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			names := *pNames
			var j *kafka.Message
			json.Unmarshal(msg.Value, j)
			if len(names) == cap(names) {
				names = names[1:]
			}
			names = append(names, "Hello")

		case <-done:
			break ConsumerLoop
		}
	}
}

func displayNames(gui *gocui.Gui, pNames *[]string) {
	for {
		gui.Execute(func(g *gocui.Gui) error {
			v, err := g.View("mainscreen")
			if err != nil {
				panic(err)
			}
			v.Clear()
			//			j, _ := json.Marshal(pNames)
			fmt.Fprintln(v, "Hi!")
			return nil
		})

		time.Sleep(2 * time.Millisecond)
	}
}
