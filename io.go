package shuttle

import (
	"bufio"
	"bytes"
	"io"
)

func Copy(ch chan<- []byte, r io.Reader, w io.Writer) error {
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := scanner.Bytes()
		cpy := make([]byte, len(line)+1)
		copy(cpy, line)
		line = append(line, '\n')

		if _, err := io.Copy(w, bytes.NewReader(line)); err != nil {
			return err
		}

		ch <- line
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
