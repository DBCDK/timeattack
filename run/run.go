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

type Request struct {
	runAt time.Time
	url   string
}

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
	chanSent := make(chan Request)
	chanDone := make(chan Request)
	chanSuccess := make(chan Request)
	chanSkipped := make(chan Request)
	chanLag := make(chan time.Duration)

	pendingRequests := make(chan Request, 100)
	requestValve := make(chan int, 100)

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

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			request, more := <-pendingRequests
			if more {
				if !*flood {
					<-requestValve
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
						chanSent <- request
						chanLag <- -request.runAt.Sub(time.Now())
						var success, _ = performRequest(request.url)
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
		url := *prefix + lineParts[1]
		var runAt time.Time

		if first {
			tCorrection = delay
			first = false
		}

		if !*flood {
			runAt = t0.Add(secsToDuration((delay - tCorrection) / *speedup))
			sleepUntil(runAt)
			requestValve <- 1
		} else {
			runAt = time.Now()
		}

		pendingRequests <- Request{runAt, url}
	}

	close(pendingRequests)
}
