package run

import (
	"bufio"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

func secsToDuration(delay float64) time.Duration {
	return time.Duration(delay * 1000000000)
}

func calcRampUpPercentage(tStart time.Time, rampUpSecs int) float64 {
	tNow := time.Now()
	rampUpDuration := time.Duration(rampUpSecs) * time.Second
	tRampUpDone := tStart.Add(rampUpDuration)
	rampUpLeft := tRampUpDone.Sub(tNow)

	res := 1 - float64(rampUpLeft)/float64(rampUpDuration)
	if res > 1 {
		return 1
	} else {
		return res
	}
}

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
	maxSleep := time.Second
	tDone := t.Sub(time.Now())
	c <- -tDone

	for tDone > 0 {
		if tDone < 10*time.Millisecond {
			time.Sleep(tDone)
		} else if tDone > maxSleep {
			time.Sleep(maxSleep)
		} else {
			time.Sleep(tDone / 2)
		}

		tDone = t.Sub(time.Now())
		c <- -tDone
	}
}

func printStatusHeader() {
	fmt.Println("runtime      lag [ms]       sent      done   waiting        successful   skipped")
}

func printStatus(tStart time.Time, sent int, done int, successful int, lag time.Duration, skipped int) {
	runTime := int(time.Now().Sub(tStart).Seconds())
	nanoLag := float64(lag.Nanoseconds()) / float64(1000000)
	fmt.Printf("\r%6ds  %12.3f   %8d  %8d  %8d  %8d %6.2f%%  %8d", runTime, nanoLag, sent, done, sent-done, successful, 100*float64(successful)/float64(done), skipped)
}

func Run(prefix *string, flood *bool, speedup *float64, rampUpSecs *int, concurrentReqs *int, requestLimit *int) {
	var wg sync.WaitGroup
	defer fmt.Println()
	defer fmt.Println()
	defer wg.Wait()

	t0 := time.Now()
	var tCorrection float64

	terminate := false
	first := true

	chanSync := make(chan int, *concurrentReqs)
	chanLimit := make(chan int, *requestLimit)
	chanSent := make(chan int)
	chanDone := make(chan int)
	chanSuccess := make(chan int)
	chanSkipped := make(chan int)
	chanLag := make(chan time.Duration)

	printStatusHeader()

	go func() {
		sent := 0
		done := 0
		success := 0
		skipped := 0
		v := 0
		var lag time.Duration

		for {
			select {
			case v = <-chanSent:
				sent += v
			case v = <-chanDone:
				done += v
				if *concurrentReqs > 0 {
					<-chanSync
				}
			case v = <-chanSuccess:
				success += v
			case v = <-chanSkipped:
				skipped += v
			case lag = <-chanLag:
			case <-time.After(100 * time.Millisecond):
				// noop ensure update of runtime
			}

			printStatus(t0, sent, done, success, lag, skipped)
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	for !terminate && scanner.Scan() {
		line := scanner.Text()

		lineParts := strings.SplitN(line, " ", 2)
		delay, _ := strconv.ParseFloat(lineParts[0], 64)
		url := *prefix + lineParts[1]

		if first {
			tCorrection = delay
			first = false
		}

		if !*flood {
			runAt := t0.Add(secsToDuration((delay - tCorrection) / *speedup))
			sleepUntil(runAt, chanLag)
		}

		wg.Add(1)
		go func() {
			defer wg.Done()

			if calcRampUpPercentage(t0, *rampUpSecs) >= rand.Float64() {
				if *requestLimit > 0 {
					select {
					case chanLimit <- 1:
					default:
						terminate = true
						return
					}
				}

				if *concurrentReqs > 0 {
					chanSync <- 1
				}
				chanSent <- 1
				var success, _ = performRequest(url)
				chanDone <- 1
				if success {
					chanSuccess <- 1
				}
			} else {
				chanSkipped <- 1
			}
		}()
	}
}
