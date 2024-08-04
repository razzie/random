package random

import (
	"errors"
	"sync"
)

var ErrFeedClosed = errors.New("feed closed")

type CancelFunc func()

// Feed is a channel that can be used to receive random uint64 values
type Feed <-chan uint64

// NewFeed returns a new Feed and its cancel function
func NewFeed(seed <-chan uint64) (Feed, CancelFunc) {
	feed := make(chan uint64)
	cancelChan := make(chan struct{})

	go func() {
		source := NewXor64Source(<-seed)
		defer close(feed)

		for {
			select {
			case <-cancelChan:
				return
			case feed <- source.Uint64():
			}

			select {
			case <-cancelChan:
				return
			case feed <- source.Uint64():
			case s := <-seed:
				source.Seed(s)
			}
		}
	}()

	return feed, sync.OnceFunc(func() { close(cancelChan) })
}

func (f Feed) Read(p []byte) (n int, err error) {
	var pos int
	var val uint64
	for n = 0; n < len(p); n++ {
		if pos == 0 {
			var ok bool
			val, ok = <-f
			if !ok {
				err = ErrFeedClosed
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
