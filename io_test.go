package shuttle

import (
	"bytes"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCopy(t *testing.T) {
	ch := make(chan []byte)

	buf := new(bytes.Buffer)
	buf.WriteString("foo\n")
	buf.WriteString("bar ")
	buf.WriteString("baz\n")

	go Copy(ch, buf)

	assert.Equal(t, "foo\n", string(<-ch))
	assert.Equal(t, "bar baz\n", string(<-ch))
}

func TestCopyLongLine(t *testing.T) {
	ch := make(chan []byte)

	buf := new(bytes.Buffer)
	buf.Grow(128 * 1024)
	for i := 0; i < 128*1024; i++ {
		buf.WriteByte('a')
	}

	go Copy(ch, buf)

	select {
	case <-time.After(time.Second * 1):
		assert.Fail(t, "Timed out expecting line")
	case line := <-ch:
		assert.Equal(t, 128*1024, len(line))
	}
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
