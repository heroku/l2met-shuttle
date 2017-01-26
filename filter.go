package shuttle

import (
	"bytes"
	"regexp"
)

func ExtractMetrics(in <-chan []byte) <-chan []byte {
	out := make(chan []byte)

	re := regexp.MustCompile(`\b((?:(?:count|measure|sample|unique)#[^=]+|source)=[^ ]+)\b`)
	go extractMetrics(re, in, out)

	return out
}

func extractMetrics(re *regexp.Regexp, in <-chan []byte, out chan<- []byte) {
	for p := range in {
		metrics := re.FindAll(p, -1)
		if len(metrics) == 0 {
			continue
		}

		out <- append(bytes.Join(metrics, []byte(" ")), []byte("\n")...)
	}

	close(out)
}
