package shuttle

import (
	"bufio"
	"bytes"
	"io"
)

func Copy(ch chan<- []byte, r io.Reader) error {
	reader := bufio.NewReader(r)

	for {
		bytes, err := reader.ReadBytes('\n')
		if len(bytes) > 0 {
			ch <- bytes
		}
		if err == io.EOF {
			return nil
		} else if err != nil {
			return err
		}
	}
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
