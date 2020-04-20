package seed

// Seeder is a function that returns a seed value or an error
type Seeder func() (uint64, error)
