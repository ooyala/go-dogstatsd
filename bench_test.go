package dogstatsd

import (
	"io"
	"io/ioutil"
	"testing"
)

type NopWriteCloser struct {
	io.Writer
}

func (*NopWriteCloser) Close() error { return nil }

func NewNopClient() *Client {
	return &Client{
		conn: &NopWriteCloser{Writer: ioutil.Discard},
	}
}

func BenchmarkGauge(b *testing.B) {
	c := NewNopClient()
	for i := 0; i < b.N; i++ {
		c.Gauge("test.gauge", 42.42, nil, 1.)
	}
}

func BenchmarkCount(b *testing.B) {
	c := NewNopClient()
	for i := 0; i < b.N; i++ {
		c.Count("test.count", 42, nil, 1.)
	}
}

func BenchmarkHistogram(b *testing.B) {
	c := NewNopClient()
	for i := 0; i < b.N; i++ {
		c.Histogram("test.histogram", 42.42, nil, 1.)
	}
}

func BenchmarkSet(b *testing.B) {
	c := NewNopClient()
	for i := 0; i < b.N; i++ {
		c.Set("test.set", "42.42", nil, 1.)
	}
}

func BenchmarkEvent(b *testing.B) {
	c := NewNopClient()
	for i := 0; i < b.N; i++ {
		c.Info("test.event.info", "Event text", nil)
	}
}

func BenchmarkSend(b *testing.B) {
	c := NewNopClient()
	for i := 0; i < b.N; i++ {
		c.send("test.send", "value", nil, 1)
	}
}
