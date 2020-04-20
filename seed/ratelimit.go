package seed

import (
	"log"
	"time"
)

// RateLimit returns a rate limited channel of seed values using the given Seeder
func RateLimit(seeder Seeder, limit int, timeframe time.Duration) <-chan uint64 {
	ch := make(chan uint64, 1)

	go func() {
		throttle := time.NewTicker(timeframe / time.Duration(limit))
		for ; true; <-throttle.C {
			seed, err := seeder.safeCall()
			if err != nil {
				log.Print("seeder error:", err)
				continue
			}

			ch <- seed
		}
	}()

	return ch
}
