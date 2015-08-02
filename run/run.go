package run

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

func Run(prefix *string, flood *bool, speedup *float64, rampUpSecs *int, concurrentReqs *int, requestLimit *int) {
	var wg sync.WaitGroup
	defer fmt.Println()
	defer wg.Wait()

	t0 := time.Now()
	var tCorrection float64

	terminate := false
	first := true

	chanSync := make(chan int, *concurrentReqs)
	chanLimit := make(chan int, *requestLimit)
	chanSent := make(chan string)
	chanDone := make(chan string)
	chanSuccess := make(chan string)
	chanSkipped := make(chan string)
	chanLag := make(chan time.Duration)

	printStatusHeader()

	go func() {
		sent := 0
		done := 0
		success := 0
		skipped := 0
		var lag time.Duration

		for {
			select {
			case _ = <-chanSent:
				sent += 1
			case _ = <-chanDone:
				done += 1
				if *concurrentReqs > 0 {
					<-chanSync
				}
			case _ = <-chanSuccess:
				success += 1
			case _ = <-chanSkipped:
				skipped += 1
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
				chanSent <- url
				var success, _ = performRequest(url)
				chanDone <- url
				if success {
					chanSuccess <- url
				}
			} else {
				chanSkipped <- url
			}
		}()
	}
}
