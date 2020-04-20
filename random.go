package random

import (
	"fmt"
)

// Feed is a channel that can be used to receive random uint64 values
type Feed <-chan uint64

// NewFeed returns a new Feed and a function that is used to close it
func NewFeed(seed <-chan uint64) (Feed, func()) {
	rf := make(chan uint64, 0)
	end := make(chan struct{}, 0)

	go func() {
		source := NewXor64Source(<-seed)
		defer close(rf)

		for {
			select {
			case <-end:
				return
			case rf <- source.Uint64():
			}

			select {
			case <-end:
				return
			case rf <- source.Uint64():
			case s := <-seed:
				source.Seed(s)
			}
		}
	}()

	return rf, func() { end <- struct{}{} }
}

func (f Feed) Read(p []byte) (n int, err error) {
	var pos int
	var val uint64
	for n = 0; n < len(p); n++ {
		if pos == 0 {
			var ok bool
			val, ok = <-f
			if !ok {
				err = fmt.Errorf("channel closed")
				return
			}
			pos = 7
		}
		p[n] = byte(val)
		val >>= 8
		pos--
	}
	return
}
