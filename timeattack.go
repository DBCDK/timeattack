package main

import (
	"github.com/dbcdk/timeattack/parse"
	"github.com/dbcdk/timeattack/run"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"runtime"
	"strconv"
)

var (
	app          = kingpin.New("timeattack", "Replays http requests").Version("1.0")
	runCmd       = kingpin.Command("run", "Start replaying requests.")
	floodCmd     = runCmd.Command("flood", "Start replaying requests")
	timedCmd     = runCmd.Command("timed", "Start replaying requests")
	tickerCmd    = runCmd.Command("ticker", "Start replaying requests")
	tickerFreq   = tickerCmd.Flag("freq", "Ticker frequency [1/s]").Default("1").Float()
	prefix       = kingpin.Flag("prefix", "Scheme, host and path prefix to prepend to urls").PlaceHolder("http://example.com").Default("").String()
	speedup      = timedCmd.Flag("speedup", "Change replay speed; 1 is 100%.").Default("1").Float()
	rampUp       = kingpin.Flag("ramp-up", "Increase the amount of requests let through over a number of seconds.").Default("0").Int()
	concurrency  = kingpin.Flag("concurrency", "Allowed concurrent requests. 0 is unlimited.").Default("0").Int()
	limit        = kingpin.Flag("limit", "Maximum number of requests that will be sent. 0 is unlimited.").Default("0").Int()
	parseCmd     = kingpin.Command("parse", "Parses input into suitable format")
	parseSolrCmd = parseCmd.Command("solr", "Parses solr logs into suitable format")
)

func init() {
	maxProcs := int64(runtime.NumCPU())
	maxProcsEnv := os.Getenv("GOMAXPROCS")
	if maxProcsEnv != "" {
		maxProcs, _ = strconv.ParseInt(maxProcsEnv, 10, 0)
	}
	runtime.GOMAXPROCS(int(maxProcs))
}

func main() {
	switch kingpin.Parse() {
	case timedCmd.FullCommand():
		run.Run(run.Timed{*speedup}, prefix, rampUp, concurrency, limit)
	case floodCmd.FullCommand():
		run.Run(run.Flood{}, prefix, rampUp, concurrency, limit)
	case tickerCmd.FullCommand():
		run.Run(run.Ticker{*tickerFreq}, prefix, rampUp, concurrency, limit)
	case parseSolrCmd.FullCommand():
		parse.Solr(prefix)
	}
}
