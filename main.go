package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	"log"
	"time"
)

type producer struct {
	// createWidget time.Time
	name    string
	setTime time.Time
}

type widget struct {
	id     string
	source string
	time   time.Time
	broken bool
}

type consumer struct {
	name string
}

// Flags for widgets, producers, consumers, number of widgets_broken
var widgetFlag int
var producerFlag int
var consumerFlag int
var widgetBrokenFlag int

// Array of Structs
var arrayOfProducers []producer
var arrayOfWidgets []widget
var arrayOfConsumers []consumer

// Command line for flagging
func commandFlag() {
	flag.IntVar(&widgetFlag, "n", 0, "an int")
	flag.IntVar(&producerFlag, "p", 0, "an int")
	flag.IntVar(&consumerFlag, "c", 0, "an int")
	flag.IntVar(&widgetBrokenFlag, "k", 0, "an int")
}

// Example UUID - 91f023ae62e42b44-af6103c770dc1dda
func generateUUID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}

	uuid := fmt.Sprintf("%x-%x",
		b[0:8], b[8:])
	return uuid
}

// func consumerConsumes() {

// 	msg := "%s consumes [%s source=%s time=%s broken=%t] in %f time"
// 	c := consumer{name: "consumer_55"}
// 	p := producer{createWidget: time.Now(), name: "producer_1"}
// 	w := widget{id: generateUUID(), source: p.name, time: p.createWidget, broken: false}
// 	fmt.Printf(msg, c.name, w.id, p.name, string(w.time.Format("15:04:05.999999")), w.broken, 22.33)
// 	fmt.Println("")
// }

// Function to produce the producers + widgets + consumers
// The algorithm every producer will be given a specific time to create the widget
func createProducer(x int) []producer {
	for i := 0; i <= x; i++ {
		producerNumber := fmt.Sprintf("producer_%d", i)
		p := producer{name: producerNumber}
		arrayOfProducers = append(arrayOfProducers, p)
	}

	return arrayOfProducers
}

// Need to use Interface
func createConsumer(x int) []consumer {
	for i := 0; i <= x; i++ {
		consumerNumber := fmt.Sprintf("consumer_%d", i)
		c := consumer{name: consumerNumber}
		arrayOfConsumers = append(arrayOfConsumers, c)
	}

	return arrayOfConsumers
}

func producerSetTime(producers []producer, c chan producer) {
	for _, p := range producers {
		go func(p producer) {
			p.setTime = time.Now()
			fmt.Println(p)
			c <- p
		}(p)
	}
}

func createWidget(c chan producer, broken bool, consumer chan widget) {
	producer := <-c
	w := widget{id: generateUUID(), source: producer.name, time: producer.setTime, broken: broken}
	// fmt.Println(w)
	consumer <- w
}

// Create a struct of Widget and Producer or for loop the Consumer struct
func consumerConsume(c chan widget, consumerNumber int) {
	// message := fmt.Printf("consumer consumes", <-c)
	consumer := arrayOfConsumers[consumerNumber]
	widget := <-c
	msg := "%s consumes [%s source=%s time=%s broken=%t] in %f time"
	fmt.Printf(msg, consumer.name, widget.id, widget.source, string(widget.time.Format("15:04:05.999999")), widget.broken, 22.33)
	fmt.Println("")
}

func main() {

	commandFlag()
	flag.Parse()
	fmt.Println("Widgets", widgetFlag)
	fmt.Println("Producers", producerFlag)
	fmt.Println("Consumers", consumerFlag)
	fmt.Println("Widgets Broken", widgetBrokenFlag)

	producerC := make(chan producer)
	// consumerC := make(chan consumer)
	widgetC := make(chan widget)

	// Blocking
	createProducer(producerFlag)
	// Blocking
	createConsumer(consumerFlag)

	// Blocking
	for _, p := range arrayOfProducers {
		fmt.Println(p.name)
	}

	// Blocking
	for _, c := range arrayOfConsumers {
		fmt.Println(c.name)
	}

	producerSetTime(arrayOfProducers, producerC)

	for {
		go createWidget(producerC, false, widgetC)
		go consumerConsume(widgetC, 0)
	}

	// Create a widget here by producer because producer has to set the time

	// d := 100 * time.Microsecond
	// fmt.Println(d)
	// fmt.Println(float64(time.Now().Nanosecond()) / 1000)
	// consumerConsumes()
	// createProducer(10)
	// fmt.Println(arrayOfProducers)
}

func chanelExample() {
	c := make(chan string)

	for i := 1; i <= 10; i++ {
		go func() {
			c <- time.Now().Format("15:04:05.999999")
		}()
	}

	for i := 1; i <= 10; i++ {
		msg := <-c
		fmt.Println(msg)
	}
}
