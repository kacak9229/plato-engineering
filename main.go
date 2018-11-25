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

// Struct based on the requirement
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

/*
1. Function to store all the flags to be input on terminal
2. Example - go run main.go -n 10 -p 1 -c 1 -k 3
*/
func commandFlag() {
	flag.IntVar(&widgetFlag, "n", 0, "an int")
	flag.IntVar(&producerFlag, "p", 0, "an int")
	flag.IntVar(&consumerFlag, "c", 0, "an int")
	flag.IntVar(&widgetBrokenFlag, "k", 0, "an int")
}

/*
1. Generate Unique universal ID for each widget
2. Example UUID - 91f023ae62e42b44-af6103c770dc1dda
*/
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

// Time difference between widget creation time and consumer consuming the widget time
func timeDifference(consumerConsume time.Time, widgetCreated time.Time) time.Duration {
	timeDiff := consumerConsume.Sub(widgetCreated)
	return timeDiff
}

// It will randomly select true or else
func randBool() bool {
	return mathRandom.Int()%2 == 0
}

// if true and the widgetBrokenFlag is greater than 0 then decrement the widgetBrokenFlag
func brokenWidget(widgetFlag int) bool {
	broken := randBool()
	if broken && (widgetBrokenFlag > 0) {
		widgetBrokenFlag--
		return true
	}
	return false
}

/*
THE ALGORITHM:
1. listen on the job channel (only receiving)
2. set the time in the setTime variable
3. use brokenWidget function to randomly grab true or false
4. create widget object and store all the value in the respective attribute
5. Finally send to both consumes channel
*/
func createWidget(id int, jobs <-chan int, consumes chan<- widget) {
	// jobs is a receving channel only
	for range jobs {
		producerID := fmt.Sprintf("producer_%d", id)
		fmt.Println("producer_", id, "setting time for widget ...")
		// time.Sleep(time.Second * 1)
		setTime := time.Now()
		IswidgetBroken := brokenWidget(widgetBrokenFlag)
		w := widget{id: generateUUID(), source: producerID, time: setTime, broken: IswidgetBroken}
		fmt.Println("producer_", id, "finished creating widget ", w.id, "at", w.time.Format("15:04:05.999999"))
		// time.Sleep(time.Second * 1)

		consumes <- w
	}
}

/*
THE ALGORITHM:
1. Check if brokenWidget is true or false
2. if broken then print out stop production message
3. else print consumer consumes message
*/
func widgetQuality(brokenWidget bool, consumerID string, widgetID string, widgetProducer string, t time.Time, timeConsume time.Duration) {
	if brokenWidget {
		fmt.Println(consumerID, "consumes [id=", widgetID, "source=", widgetProducer, "time=", t.Format("15:04:05.999999"), "broken=", brokenWidget, "-- stopping production [execution stops]")
		// time.Sleep(time.Second * 1)
		os.Exit(0)
	}
	fmt.Println(consumerID, "consumes [id=", widgetID, "source=", widgetProducer, "time=", t.Format("15:04:05.999999"), "broken=", brokenWidget, "in", timeConsume, "time")
}

/*
THE ALGORITHM:
1. Listen on consumes channel
2. create the time consumer consumes
3. Use timeDifference function to deduct two different time
4. call widgetQuality function to see a different message based on true or false of the widget
5. finally send to done channel
*/
func consumerConsume(id int, done chan<- widget, consumes <-chan widget) {
	for w := range consumes {
		consumerID := fmt.Sprintf("consumer_%d", id)
		consumerConsumeTime := time.Now()
		timeConsume := timeDifference(consumerConsumeTime, w.time)
		widgetQuality(w.broken, consumerID, w.id, w.source, w.time, timeConsume)
		done <- w
	}
}

/* Main function */
func main() {
	commandFlag()
	flag.Parse()

	// Create three seperate channels, 1 is for job
	jobs := make(chan int, 50)
	consumes := make(chan widget, 50)
	done := make(chan widget, 50)

	// Loop the producer flag and create a goroutine for widget creation
	for w := 1; w <= producerFlag; w++ {
		go createWidget(w, jobs, consumes)
	}

	// Send to job
	for j := 1; j <= widgetFlag; j++ {
		jobs <- j
	}

	// Loop the producer flag and create a goroutine for consumer consume the widget
	for c := 1; c <= consumerFlag; c++ {
		go consumerConsume(c, done, consumes)
	}

	// Done
	for a := 1; a <= widgetFlag; a++ {
		<-done
	}

}
