package main

import (
	"github.com/dbcdk/timeattack/parse"
	"github.com/dbcdk/timeattack/run"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app          = kingpin.New("timeattack", "Replays http requests").Version("1.0")
	runCmd       = kingpin.Command("run", "Start replaying requests")
	prefix       = kingpin.Flag("prefix", "Scheme, host and path prefix to prepend to urls").PlaceHolder("http://example.com").Default("").String()
	flood        = runCmd.Flag("flood", "Run requests as fast as possible, ignoring timings").Bool()
	speedup      = runCmd.Flag("speedup", "Manipulate time between requests. Values above '1' increase the speed, while values below decreases speed.").Default("1").Float()
	parseCmd     = kingpin.Command("parse", "Parses input into suitable format")
	parseSolrCmd = parseCmd.Command("solr", "Parses solr logs into suitable format")
)

func main() {
	switch kingpin.Parse() {
	case runCmd.FullCommand():
		run.Run(prefix, flood, speedup)
	case parseSolrCmd.FullCommand():
		parse.Solr(prefix)
	}
}