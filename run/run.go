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

func Run(scheduler Scheduler, prefix *string, rampUpSecs *int, concurrentReqs *int, requestLimit *int) {
	var wg sync.WaitGroup
	defer fmt.Println()
	defer wg.Wait()

	t0 := time.Now()
	var tCorrection float64

	terminate := false
	first := true

	chanSync := make(chan int, *concurrentReqs)
	chanLimit := make(chan int, *requestLimit)
	chanSent := make(chan Request)
	chanDone := make(chan Request)
	chanSuccess := make(chan Request)
	chanSkipped := make(chan Request)
	chanLag := make(chan time.Duration)

	pendingRequests := make(chan Request, 100)

	schedulerInput := make(chan Request)
	defer close(schedulerInput)
	go scheduler.Run(schedulerInput, pendingRequests)

	printStatusHeader()

	go func() {
		ticks := time.Tick(100 * time.Millisecond)
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
			case <-ticks:
				// noop ensure update of runtime
			}

			printStatus(t0, sent, done, success, lag, skipped)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			request, more := <-pendingRequests
			if more {
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
						chanSent <- request
						chanLag <- -request.tStart.Add(secsToDuration(request.delay)).Sub(time.Now())
						var success = performRequest(request.url)
						chanDone <- request
						if success {
							chanSuccess <- request
						}
					} else {
						chanSkipped <- request
					}
				}()
			} else {
				return
			}
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	for !terminate && scanner.Scan() {
		line := scanner.Text()

		lineParts := strings.SplitN(line, " ", 2)
		delay, _ := strconv.ParseFloat(lineParts[0], 64)

		if first {
			tCorrection = delay
			first = false
		}

		request := Request{t0, delay - tCorrection, *prefix + lineParts[1]}

		schedulerInput <- request
	}
}
