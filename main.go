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

type widget struct {
	id     string
	source string
	time   time.Time
	broken bool
}

// Flags for widgets, producers, consumers, number of widgets_broken
var widgetFlag int
var producerFlag int
var consumerFlag int
var widgetBrokenFlag int

var widgets map[int]widget

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

// Time difference
func timeDifference(widgetCreated time.Time, consumerConsume time.Time) time.Duration {
	diff := widgetCreated.Sub(consumerConsume)
	return diff
}

// It will randomly select true or else
func randBool() bool {
	return mathRandom.Int()%2 == 0
}

func brokenWidget(widgetFlag int) bool {
	broken := randBool()
	if broken && (widgetBrokenFlag > 0) {
		widgetBrokenFlag -= 1
		return true
	}
	return false
}

func createWidget(id int, jobs <-chan int, results chan<- int, consumes chan<- widget) {
	for j := range jobs {
		producer_id := fmt.Sprintf("producer_%d", id)
		fmt.Println("producer_", id, "setting time for widget ...")
		setTime := time.Now()
		time.Sleep(time.Second * 1)
		widgetBreaks := brokenWidget(widgetBrokenFlag)
		w := widget{id: generateUUID(), source: producer_id, time: setTime, broken: widgetBreaks}
		fmt.Println("producer_", id, "finished creating widget ", w.id, "at", w.time.Format("15:04:05.999999"))
		time.Sleep(time.Second * 1)
		results <- j
		consumes <- w
	}
}

func widgetQuality(brokenWidget bool, consumer_id string, widgetId string, widgetProducer string, t time.Time, timeConsume time.Duration) {
	if brokenWidget {
		fmt.Println(consumer_id, "consumes [id=", widgetId, "source=", widgetProducer, "time=", t.Format("15:04:05.999999"), "broken=", brokenWidget, "-- stopping production [execution stops]")
		time.Sleep(time.Second * 2)
		os.Exit(0)
	}
	fmt.Println(consumer_id, "consumes [id=", widgetId, "source=", widgetProducer, "time=", t.Format("15:04:05.999999"), "broken=", brokenWidget, "in", timeConsume, "time")
}

func consumerConsume(id int, results <-chan int, done chan<- widget, consumes <-chan widget) {
	for w := range consumes {
		consumer_id := fmt.Sprintf("consumer_%d", id)
		consumerConsume := time.Now()
		timeConsume := timeDifference(consumerConsume, w.time)
		widgetQuality(w.broken, consumer_id, w.id, w.source, w.time, timeConsume)
		done <- w
	}
}

func main() {
	commandFlag()
	flag.Parse()

	jobs := make(chan int, 100)
	results := make(chan int, 100)
	consumes := make(chan widget, 100)
	done := make(chan widget, 100)

	for w := 1; w <= producerFlag; w++ {
		go createWidget(w, jobs, results, consumes)
	}

	for j := 1; j <= widgetFlag; j++ {
		jobs <- j
	}
	// close(jobs)

	for c := 1; c <= consumerFlag; c++ {
		go consumerConsume(c, results, done, consumes)
	}

	for a := 1; a <= widgetFlag; a++ {
		<-done
	}

}
