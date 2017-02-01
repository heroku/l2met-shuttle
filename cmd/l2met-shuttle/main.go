package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/heroku/l2met-shuttle"
	lshuttle "github.com/heroku/log-shuttle"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: %v <url>\n", os.Args[0])
		os.Exit(1)
	}

	url := os.Args[1]

	ch := make(chan []byte)

	config := lshuttle.NewConfig()
	config.LogsURL = url

	r := shuttle.Reader(shuttle.SkipLuhnMatches(shuttle.ExtractMetrics(ch)))

	s := lshuttle.NewShuttle(config)
	s.LoadReader(ioutil.NopCloser(r))
	s.Launch()

	shuttle.Copy(ch, os.Stdin)

	close(ch)

	s.WaitForReadersToFinish()
	s.Land()
}
