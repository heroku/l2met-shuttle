package shuttle

import (
	"bufio"
	"bytes"
	"io"
)

func Copy(ch chan<- []byte, r io.Reader, w io.Writer) error {
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := scanner.Text() + "\n"
		if err := writeFully(w, line); err != nil {
			return err
		}

		ch <- []byte(line)
	}

	return scanner.Err()
}

func writeFully(w io.Writer, str string) error {
	bytes := []byte(str)
	offset := 0

	for offset < len(bytes) {
		slice := bytes[offset:]

		written, err := w.Write(slice)
		if err != nil {
			return err
		}

		offset += written
	}

	return nil
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
