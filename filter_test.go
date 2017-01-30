package shuttle

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractMetrics(t *testing.T) {
	ch := make(chan []byte)
	r := ExtractMetrics(ch)

	go func() {
		ch <- []byte("foo count#requests=1\n")
		ch <- []byte("measure#query=0.2ms bar baz\n")
		ch <- []byte("qux\n")
		ch <- []byte("quux sample#size=0 unique#user=alice\n")
		ch <- []byte("source=development\n")
	}()

	assert.Equal(t, "count#requests=1\n", string(<-r))
	assert.Equal(t, "measure#query=0.2ms\n", string(<-r))
	assert.Equal(t, "sample#size=0 unique#user=alice\n", string(<-r))
	assert.Equal(t, "source=development\n", string(<-r))
}

func TestCloseExtractMetrics(t *testing.T) {
	ch := make(chan []byte)
	r := ExtractMetrics(ch)

	go close(ch)

	_, ok := <-r
	assert.False(t, ok)
}

func TestSkipLuhnMatches(t *testing.T) {
	ch := make(chan []byte, 5)
	ch <- []byte("4111 1111 1111 1111\n")
	ch <- []byte("measure#latency=4.111111111111111\n")
	ch <- []byte("sample#depth=4222222222222222\n")
	ch <- []byte("unique#card=4111111111111111\n")
	ch <- []byte("\n")

	r := SkipLuhnMatches(ch)

	assert.Equal(t, "sample#depth=4222222222222222\n", string(<-r))
	assert.Equal(t, "\n", string(<-r))
}
