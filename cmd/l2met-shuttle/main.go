package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/heroku/l2met-shuttle"
	lshuttle "github.com/heroku/log-shuttle"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %v [options] <url>\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(1)
}

func parseArgs() (string, io.Writer) {
	tee := flag.Bool("tee", false, "pipe input through to stdout")
	flag.Parse()

	if flag.NArg() != 1 {
		usage()
	}

	var out io.Writer = ioutil.Discard
	if *tee {
		out = os.Stdout
	}

	return flag.Arg(0), out
}

func main() {
	url, output := parseArgs()

	ch := make(chan []byte)

	config := lshuttle.NewConfig()
	config.LogsURL = url

	r := shuttle.Reader(shuttle.SkipLuhnMatches(shuttle.ExtractMetrics(ch)))

	s := lshuttle.NewShuttle(config)
	s.LoadReader(ioutil.NopCloser(r))
	s.Launch()

	shuttle.Copy(ch, io.TeeReader(os.Stdin, output))

	close(ch)

	s.WaitForReadersToFinish()
	s.Land()
}
