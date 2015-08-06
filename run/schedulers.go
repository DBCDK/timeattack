package run

import (
	"time"
)

type Scheduler interface {
	Run(input chan Request, output chan Request)
}

type Timed struct {
	Speedup float64
}

func (timed Timed) Run(input chan Request, output chan Request) {
	for {
		request, ok := <-input
		request.delay = request.delay / timed.Speedup
		if ok {
			sleepUntil(request.tStart.Add(secsToDuration(request.delay)))
			output <- request
		} else {
			close(output)
			return
		}
	}
}

type Flood struct{}

func (flood Flood) Run(input chan Request, output chan Request) {
	for {
		request, ok := <-input
		request.delay = 0
		if ok {
			output <- request
		} else {
			close(output)
			return
		}
	}
}

type Ticker struct {
	Frequency float64 // [1/s]
}

func (ticker Ticker) Run(input chan Request, output chan Request) {
	c := time.Tick(secsToDuration(1 / float64(ticker.Frequency)))
	requests := 0
	for range c {
		request, ok := <-input
		requests += 1
		request.delay = float64(requests) / float64(ticker.Frequency)
		if ok {
			output <- request
		} else {
			close(output)
			return
		}
	}
}
