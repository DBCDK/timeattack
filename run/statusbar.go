package run

import (
	"fmt"
	"time"
)

func printStatusHeader() {
	fmt.Println("runtime      lag [ms]       sent      done   waiting        successful   skipped")
}

func printStatus(tStart time.Time, sent int, done int, successful int, lag time.Duration, skipped int) {
	runTime := int(time.Now().Sub(tStart).Seconds())
	nanoLag := float64(lag.Nanoseconds()) / float64(1000000)
	fmt.Printf("\r%6ds  %12.3f   %8d  %8d  %8d  %8d %6.2f%%  %8d", runTime, nanoLag, sent, done, sent-done, successful, 100*float64(successful)/float64(done), skipped)
}
