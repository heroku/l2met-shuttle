package shuttle

import (
	"bytes"
	"testing"

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
