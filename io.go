package shuttle

import (
	"bufio"
	"bytes"
	"io"
)

func Copy(ch chan<- []byte, r io.Reader) error {
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		ch <- []byte(scanner.Text() + "\n")
	}

	return scanner.Err()
}

type reader struct {
	buf *bytes.Buffer
	ch  <-chan []byte
}

func Reader(ch <-chan []byte) io.Reader {
	return &reader{
		new(bytes.Buffer),
		ch,
	}
}

func (r *reader) Read(p []byte) (int, error) {
	if r.buf.Len() == 0 {
		p, ok := <-r.ch
		if !ok {
			return 0, io.EOF
		}

		_, err := r.buf.Write(p)
		if err != nil {
			return 0, err
		}
	}

	return r.buf.Read(p)
}
