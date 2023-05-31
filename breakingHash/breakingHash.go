package breakingHash

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/abramtrinh/wc3-save-code/lookup2"
)

const letters = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	// Ceil(Log2(len(letters))) == Ceil(Log2(62)) == 6
	// 6 bits can represent the 62 indices of letters
	lettersBits = 6
	// 1-bit mask (e.g. 111111) to extract indices
	lettersMask = 1<<lettersBits - 1
	// # of letters you can extract from 63 bits
	lettersMax = 63 / lettersBits
)

// Generate random strings by assigning characters based on "letters"'s indices.
func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// Generates random strings by masking a random 63-bit buffer to generate indicies for the letters.
// Should be faster than randStringBytes due to const & not having to call rand.Intn as often.
// Check benchmarks.
func randStringMask(n int) string {
	/*
		Extract 6 bits from a random 63 bit int by masking out the last 6.
		Then shifting >> to remove the recently masked bits.
		All the while keeping track of the # of extractable letters/bits left of the 63 bit int.
		When no longer extractable, a new 63 bit int is randomized.
		Note: indicies of letters is 0->61 and 6bit numbers have a max of 64. So discard 63/64.
	*/

	b := make([]byte, n)

	buffer := rand.Int63()
	remaining := lettersMax

	for i := 0; i < n; {
		if remaining == 0 {
			buffer = rand.Int63()
			remaining = lettersMax
		}

		maskedValue := int(buffer & lettersMask)
		if maskedValue < len(letters) {
			b[i] = letters[maskedValue]
			i++
		}

		buffer = buffer >> lettersBits
		remaining--
	}

	return string(b)
}

// Spins up # of goroutines trying to find a random string that hashes to specified hash.
// Brute force approach of a first preimage attack.
func BruteForceLookup2(hash uint32, jobs int) string {
	startTime := time.Now()
	defer func() {
		fmt.Printf("Time elapsed %v \n", time.Since(startTime))
	}()

	result := make(chan string)

	for i := 0; i < jobs; i++ {
		go func() {
			for {
				// Note: Should I change the length of random string to something else?
				key := randStringMask(10)
				if hash == lookup2.Hash(key, 0) {
					result <- key
				}
			}
		}()
	}
	// Channel is used because I just want the first goroutine that finds a key to return and quit.
	return <-result
}
