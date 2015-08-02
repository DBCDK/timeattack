package run

import (
	"net/http"
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
