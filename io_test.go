package shuttle

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCopy(t *testing.T) {
	ch := make(chan []byte)

	buf := new(bytes.Buffer)
	buf.WriteString("foo\n")
	buf.WriteString("bar ")
	buf.WriteString("baz\n")

	out := new(bytes.Buffer)

	go Copy(ch, buf, out)

	assert.Equal(t, "foo\n", string(<-ch))
	assert.Equal(t, "bar baz\n", string(<-ch))

	assert.Equal(t, "foo\nbar baz\n", out.String())
}

func TestRead(t *testing.T) {
	ch := make(chan []byte, 2)
	ch <- []byte("hello")
	ch <- []byte("world")

	r := Reader(ch)

	p := make([]byte, 4)
	n, _ := r.Read(p)
	assert.Equal(t, "hell", string(p[:n]))
	n, _ = r.Read(p)
	assert.Equal(t, "o", string(p[:n]))
	p = make([]byte, 6)
	n, _ = r.Read(p)
	assert.Equal(t, "world", string(p[:n]))
}

func TestCloseRead(t *testing.T) {
	ch := make(chan []byte)
	close(ch)

	r := Reader(ch)

	p := make([]byte, 1)
	n, err := r.Read(p)
	assert.Equal(t, 0, n)
	assert.Equal(t, io.EOF, err)
}
