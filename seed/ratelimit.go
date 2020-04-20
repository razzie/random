package seed

import (
	"time"
)

// RateLimit returns a rate limited channel of seed values using the given Seeder
func RateLimit(seeder Seeder, limit int, timeframe time.Duration) <-chan uint64 {
	ch := make(chan uint64, 1)

	go func() {
		ticker := time.NewTicker(timeframe)
		for ; true; <-ticker.C {
			for i := 0; i < limit; i++ {
				seed, err := seeder()
				if err != nil {
					continue
				}

				ch <- seed
			}
		}
	}()

	return ch
}
