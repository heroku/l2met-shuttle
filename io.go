package shuttle

import (
	"bufio"
	"io"
)

func Copy(ch chan<- []byte, r io.Reader) error {
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		ch <- []byte(scanner.Text() + "\n")
	}

	return scanner.Err()
}
