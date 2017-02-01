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

func SkipLuhnMatches(in <-chan []byte) <-chan []byte {
	out := make(chan []byte)

	go skipLuhnMatches(in, out)

	return out
}

func skipLuhnMatches(in <-chan []byte, out chan<- []byte) {
	re := regexp.MustCompile(`\D+`)

	for p := range in {
		numbers := re.Split(string(p), -1)
		if containsLuhnMatch(numbers) {
			continue
		}

		out <- p
	}

	close(out)
}

func containsLuhnMatch(numbers []string) bool {
	for i := 0; i < len(numbers); i++ {
		value := ""
		for j := i; j < len(numbers); j++ {
			value += numbers[j]
			if len(value) < 13 {
				continue
			}
			if len(value) > 19 {
				break
			}
			if luhnMatch(value) {
				return true
			}
		}
	}

	return false
}

func luhnMatch(n string) bool {
	var sum int

	for alternate, i := false, len(n)-1; i > -1; i-- {
		mod := int(n[i]) - 48

		if alternate {
			mod *= 2
			if mod > 9 {
				mod = (mod % 10) + 1
			}
		}

		alternate = !alternate
		sum += mod
	}

	return sum%10 == 0
}
