package seed

import (
	"fmt"
)

// Seeder is a function that returns a seed value or an error
type Seeder func() (uint64, error)

func (seeder Seeder) safeCall() (v uint64, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	return seeder()
}
