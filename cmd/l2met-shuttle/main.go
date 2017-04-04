package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/heroku/l2met-shuttle"
	lshuttle "github.com/heroku/log-shuttle"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %v [--tee] <url>\n", os.Args[0])
	os.Exit(1)
}

func parseArgs() (string, io.Writer) {
	switch len(os.Args) {
	case 2: // [l2met-shuttle, url]
		return os.Args[1], ioutil.Discard
	case 3: // [l2met-shuttle, --tee?, url]
		if os.Args[1] == "--tee" {
			return os.Args[2], os.Stdout
		}
	}

	usage()
	return "", nil // unreachable
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

	shuttle.Copy(ch, os.Stdin, output)

	close(ch)

	s.WaitForReadersToFinish()
	s.Land()
}
