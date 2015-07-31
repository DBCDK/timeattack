package parse

import (
	"bufio"
	"fmt"
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
		timeEnd := strings.Index(line, ",")
		t, _ := time.Parse(timeLayout, line[:timeEnd])

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
			fmt.Printf("%f %s", float64(t.UnixNano())/1000000000, *prefix+"/select?"+query)
			fmt.Println()
		}
	}
}
