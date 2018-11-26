package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	"log"
	mathRandom "math/rand"
	"os"
	"time"
)

type producer struct {
	id string
}

type consumer struct {
	id string
}

type widget struct {
	id     string
	source string
	time   time.Time
	broken bool
}

var widgetFlag int
var producerFlag int
var consumerFlag int
var widgetBrokenFlag int

func (w *widget) isWidgetBroken(widgetFlag int) bool {
	broken := mathRandom.Int()%2 == 0
	if broken && (widgetFlag > 0) {
		widgetFlag--
		return true
	}
	return false
}

func (w *widget) generateUUID() {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}

	uuid := fmt.Sprintf("%x-%x",
		b[0:8], b[8:])
	w.id = uuid
}

func hireProducer(number int) producer {
	id := fmt.Sprintf("producer_%d", number)
	p := producer{id: id}
	return p
}

func (p *producer) setTime() time.Time {
	currentTime := time.Now()
	return currentTime
}

func (p *producer) createWidget() widget {
	w := widget{}
	w.generateUUID()
	w.source = p.id
	w.time = p.setTime()
	w.broken = w.isWidgetBroken(widgetBrokenFlag)
	return w
}

func consumerEnterTheShop(number int) consumer {
	id := fmt.Sprintf("consumer_%d", number)
	c := consumer{id: id}
	return c
}

func (c *consumer) consumeTime(w time.Time) time.Duration {
	consumeTime := time.Now().Sub(w)
	return consumeTime
}

func (c *consumer) consumeWidget(w widget) {
	if w.broken {
		fmt.Println(c.id, "consumes [id=", w.id, "source=", w.source,
			"time=", w.time.Format("15:04:05.999999"), "broken=", w.broken,
			"-- stopping production [execution stops]")
		os.Exit(0)
	} else {
		fmt.Println(c.id, "consumes [id=", w.id, "source=", w.source,
			"time=", w.time.Format("15:04:05.999999"), "broken=", w.broken,
			"in", c.consumeTime(w.time), "time")
	}
}

func startProduction(id int, production <-chan int, shop chan<- widget) {
	for range production {
		p := hireProducer(id)
		fmt.Println(p.id, "is setting up time for the widget ....")
		time.Sleep(time.Microsecond * 1)
		fmt.Println(p.id, "is creating the widget ....")
		time.Sleep(time.Second * 1)
		w := p.createWidget()
		fmt.Println(p.id, "finished creating widget", w.id)
		shop <- w
	}
}

func consumerBuying(id int, shop <-chan widget, done chan<- widget) {
	for w := range shop {
		consumer := consumerEnterTheShop(id)
		consumer.consumeWidget(w)
		done <- w
	}
}

func main() {

	flag.IntVar(&widgetFlag, "n", 0, "an int")
	flag.IntVar(&producerFlag, "p", 0, "an int")
	flag.IntVar(&consumerFlag, "c", 0, "an int")
	flag.IntVar(&widgetBrokenFlag, "k", 0, "an int")
	flag.Parse()

	production := make(chan int, 40)
	shop := make(chan widget, 40)
	done := make(chan widget, 40)

	for s := 1; s <= producerFlag; s++ {
		go startProduction(s, production, shop)
	}

	// Send to production channel
	for p := 1; p <= widgetFlag; p++ {
		production <- p
	}

	for c := 1; c <= consumerFlag; c++ {
		go consumerBuying(c, shop, done)
	}

	// Done
	for a := 1; a <= widgetFlag; a++ {
		<-done
	}

}
