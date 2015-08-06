package run

import (
	"io/ioutil"
	"net/http"
)

func performRequest(url string) bool {
	res, errGet := http.Get(url)

	if errGet != nil {
		// log error
		return false
	} else {
		_, errRead := ioutil.ReadAll(res.Body)
		if errRead == nil {
			res.Body.Close()

			if res.StatusCode != 200 {
				// log error
			}
		} else {
			// log error
		}

		return errRead == nil && res.StatusCode == 200
	}
}
