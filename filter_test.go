package shuttle

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractMetrics(t *testing.T) {
	ch := make(chan []byte, 5)
	ch <- []byte("foo count#requests=1 count#requests=-1\n")
	ch <- []byte("measure#query=0.2ms bar baz\n")
	ch <- []byte("qux\n")
	ch <- []byte("quux sample#size=0 unique#user=alice\n")
	ch <- []byte("source=development\n")

	r := ExtractMetrics(ch)

	assert.Equal(t, "count#requests=1 count#requests=-1\n", string(<-r))
	assert.Equal(t, "measure#query=0.2ms\n", string(<-r))
	assert.Equal(t, "sample#size=0 unique#user=alice\n", string(<-r))
	assert.Equal(t, "source=development\n", string(<-r))
}

func TestCloseExtractMetrics(t *testing.T) {
	ch := make(chan []byte)
	close(ch)

	r := ExtractMetrics(ch)

	_, ok := <-r
	assert.False(t, ok)
}
