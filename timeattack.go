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
	runCmd       = kingpin.Command("run", "Start replaying requests")
	prefix       = kingpin.Flag("prefix", "Scheme, host and path prefix to prepend to urls").PlaceHolder("http://example.com").Default("").String()
	flood        = runCmd.Flag("flood", "Run requests as fast as possible, ignoring timings").Bool()
	speedup      = runCmd.Flag("speedup", "Manipulate time between requests. Values above '1' increase the speed, while values below decreases speed.").Default("1").Float()
	rampUp       = runCmd.Flag("ramp-up", "Increase the amount of requests let through over a number of seconds.").Default("0").Int()
	concurrency  = runCmd.Flag("concurrency", "Allowed concurrent requests. 0 is unlimited.").Default("0").Int()
	limit        = runCmd.Flag("limit", "Maximum number of requests that will be sent. 0 is unlimited.").Default("0").Int()
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
	case runCmd.FullCommand():
		run.Run(prefix, flood, speedup, rampUp, concurrency, limit)
	case parseSolrCmd.FullCommand():
		parse.Solr(prefix)
	}
}
