package run

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

func performRequest(url string) (bool, time.Duration) {
	t0 := time.Now()
	res, err := http.Get(url)
	t1 := time.Now()
	if err == nil {
		res.Body.Close()
	}

	return err == nil && res.StatusCode == 200, t1.Sub(t0)
}

func sleepUntil(t time.Time, c chan time.Duration) {
	tDone := t.Sub(time.Now())
	c <- -tDone

	for tDone > 0 {
		var sleepDuration time.Duration
		if tDone < 10*time.Millisecond {
			sleepDuration = tDone
		} else {
			sleepDuration = tDone / 2
		}

		time.Sleep(sleepDuration)
		tDone = t.Sub(time.Now())
		c <- -tDone
	}
}

func printStatus(sent int, done int, successful int, lag time.Duration) {
	nanoLag := float64(lag.Nanoseconds()) / float64(1000000)
	fmt.Printf("\rsent: %8d, done: %8d, waiting: %8d, successful: %8d (%3.2f%%), lag: %8.3fms", sent, done, sent-done, successful, 100*float64(successful)/float64(done), nanoLag)
}

func Run(prefix *string, flood *bool, speedup *float64) {
	var wg sync.WaitGroup

	t0 := time.Now()
	var tCorrection float64

	first := true

	chanSent := make(chan int)
	chanDone := make(chan int)
	chanSuccess := make(chan int)
	chanLag := make(chan time.Duration)

	go func() {
		sent := 0
		done := 0
		success := 0
		v := 0
		var lag time.Duration

		for {
			select {
			case v = <-chanSent:
				sent += v
			case v = <-chanDone:
				done += v
			case v = <-chanSuccess:
				success += v
			case lag = <-chanLag:
			}

			printStatus(sent, done, success, lag)
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()

		lineParts := strings.SplitN(line, " ", 2)
		delay, _ := strconv.ParseFloat(lineParts[0], 64)
		url := *prefix + lineParts[1]

		if first {
			tCorrection = delay
			first = false
		}

		if !*flood {
			runAt := t0.Add(time.Duration((delay - tCorrection) / *speedup * 1000000000))
			sleepUntil(runAt, chanLag)
		}

		wg.Add(1)
		go func() {
			defer wg.Done()

			chanSent <- 1
			var success, _ = performRequest(url)
			chanDone <- 1
			if success {
				chanSuccess <- 1
			}
		}()
	}

	wg.Wait()
	fmt.Println()
}
