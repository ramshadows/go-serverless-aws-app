package utils

import (
	"math/rand"
	"time"
)

// This is a special function that will be called exactly once before
// any other code in the package is executed
func init() {
	// tell rand to use the current unix nano as the seed value.
	rand.Seed(time.Now().UnixNano())
}

// randomInt function uses the rand.Int() to generate from
// zero to (max-min). If we add min to it, we will get a value from min to max
// randomInt() function can be used to set the number of cores and the number of threads
// Cores would be between 2 cores and 8 cores
// The number of threads will be a random integer between the number of cores and 12
func RandomInt(min, max int) int {
	return min + rand.Int()%(max-min+1)
}
