package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/Shopify/sarama"
	"github.com/jroimartin/gocui"
	"time"
	"github.com/RobotDisco/kafkaDemo"
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
	refresh := make(chan bool)

	go streamKafka(gui, partitionConsumer, done)
	go updateNames(gui, partitionConsumer, refresh, done)
	go updateSums(gui, partitionConsumer, refresh, done)
	go timer(refresh)

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

	v, err = gui.SetView("mainscreen", 0, 10, x, 12)
	if err != nil && err != gocui.ErrUnknownView {
		panic(err)
	}

	v, err = gui.SetView("sumscreen", 0, 12, x, y)
	if err != nil && err != gocui.ErrUnknownView {
		panic(err)
	}

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

func updateNames(gui *gocui.Gui, partitionConsumer sarama.PartitionConsumer, refresh <-chan bool, done <-chan bool) {
	names := [10]string{}
	pos := 0

ConsumerLoop:
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			var j kafkaDemo.Message
			if err := json.Unmarshal(msg.Value, &j); err != nil {
				panic(nil)
			}
			//if len(names) == cap(names) {
			//	names = names[1:]
			//}
			names[pos] = j.Msg
			pos = (pos + 1) % 10
		case <-done:
			break ConsumerLoop
		case <-refresh:
			gui.Execute(func(g *gocui.Gui) error {
				v, err := g.View("mainscreen")
				if err != nil {
					panic(err)
				}
				v.Clear()
				fmt.Fprintln(v, names)
				g.SetCurrentView("mainscreen")
				return nil
			})
		}
	}
}

func updateSums(gui *gocui.Gui, partitionConsumer sarama.PartitionConsumer, refresh <-chan bool, done <-chan bool) {
	sums := make(map[string]int)
	sums[kafkaDemo.SrcIOS] = 0
	sums[kafkaDemo.SrcBrowser] = 0
	sums[kafkaDemo.SrcWeb] = 0

ConsumerLoop:
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			var j kafkaDemo.Message
			if err := json.Unmarshal(msg.Value, &j); err != nil {
				panic(nil)
			}
			sums[j.Type] = sums[j.Type] + 1

		case <-done:
			break ConsumerLoop
		case <-refresh:
			gui.Execute(func(g *gocui.Gui) error {
				v, err := g.View("sumscreen")
				if err != nil {
					panic(err)
				}
				v.Clear()
				fmt.Fprintln(v, kafkaDemo.SrcIOS+": "+strconv.Itoa(sums[kafkaDemo.SrcIOS]))
				fmt.Fprintln(v, kafkaDemo.SrcBrowser+": "+strconv.Itoa(sums[kafkaDemo.SrcBrowser]))
				fmt.Fprintln(v, kafkaDemo.SrcWeb+": "+strconv.Itoa(sums[kafkaDemo.SrcWeb]))
				g.SetCurrentView("sumscreen")
				return nil
			})
		}
	}
}

func timer(refresh chan<- bool) {
	for {
		time.Sleep(1000 * time.Millisecond)
		refresh <- true
	}
}
