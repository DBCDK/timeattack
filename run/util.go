package run

import (
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

func sleepUntil(t time.Time) {
	maxSleep := time.Second
	tDone := t.Sub(time.Now())

	for tDone > 0 {
		if tDone < 10*time.Millisecond {
			time.Sleep(tDone)
		} else if tDone > maxSleep {
			time.Sleep(maxSleep)
		} else {
			time.Sleep(tDone / 2)
		}

		tDone = t.Sub(time.Now())
	}
}
