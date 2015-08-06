package parse

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"
)

func Solr(prefix *string) {
	timeLayout := "2006-01-02 15:04:05.000"

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		isSelect := false
		var query string
		parts := strings.Split(line, " ")
		splitByComma := strings.Split(line, ",")

		if !(len(splitByComma) > 2) {
			continue
		}

		if !(splitByComma[1] == "INFO") {
			continue
		}

		t, tErr := time.Parse(timeLayout, splitByComma[0])

		if tErr != nil {
			continue
		}

		for _, part := range parts {
			if part == "path=/select" {
				isSelect = true
			}

			if strings.HasPrefix(part, "params=") {
				qStart := strings.Index(part, "{")
				qEnd := strings.Index(part, "}")
				query = part[qStart+1 : qEnd]
			}
		}

		if isSelect {
			resultingUrl := *prefix + "/select?" + query
			_, err := url.Parse(resultingUrl)
			if err != nil {
				continue
			}
			fmt.Printf("%f %s", float64(t.UnixNano())/1000000000, resultingUrl)
			fmt.Println()
		}
	}
}
