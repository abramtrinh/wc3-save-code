package breakingHash

import (
	"errors"
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
func BruteForceLookup2(hash uint32, jobs int) (string, error) {
	//Errors if goroutines is set to <= 0
	if jobs <= 0 {
		return "", errors.New("Jobs/goroutines cannot be 0 or less.")
	}

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
	return <-result, nil
}

// Reverses the lookup2.c hashing. Returns a byte slice that can be convereted to string later.
func UnHash(hash uint32) []byte {

	// You can set arbitrary value to anything.
	// length is set to 12 for easier masking and unhashing.
	// As usual, if you know hash seed changed, then you need to change seed.
	var arbitrary, length, seed uint32 = 0, 12, 0

	// Values of a, b, c are to mimic internal state of Hash.
	var a, b, c uint32 = 0x9e3779b9, 0x9e3779b9, uint32(seed)

	//First Unmix
	unA, unB, unC := unMix(arbitrary, arbitrary, hash)
	unC -= length

	//Second Unmix
	unA, unB, unC = unMix(unA, unB, unC)

	// Reverses the hashing and extract the bytes from the resulting value.
	var finalByteSlice []byte
	finalByteSlice = append(finalByteSlice, byteMasking(unA-a)...)
	finalByteSlice = append(finalByteSlice, byteMasking(unB-b)...)
	finalByteSlice = append(finalByteSlice, byteMasking(unC-c)...)

	// Note that the bytes that are in here are UTF-8 and may contain values that can't be printed.
	// e.g. UTF-8 144 
	return finalByteSlice
}

// Reverses the internal state mixing from lookup2.c
func unMix(a, b, c uint32) (uint32, uint32, uint32) {
	shiftValues := [9]uint32{15, 10, 3, 5, 16, 12, 13, 8, 13}

	for i := 0; i < len(shiftValues); i += 3 {

		// Shifts 15, 5, 13
		c = (c ^ (b >> shiftValues[i])) + a + b

		// Shifts 10, 16, 8
		b = (b ^ (a << shiftValues[i+1])) + c + a

		// Shifts 3, 12, 13
		a = (a ^ (c >> shiftValues[i+2])) + b + c

	}

	return a, b, c
}

// Splits a 32-bit number into 4 seperate bytes.
func byteMasking(num uint32) []byte {

	index0 := uint8(num & 0xFF)
	index8 := uint8(num & 0xFF00 >> 8)
	index16 := uint8(num & 0xFF0000 >> 16)
	index24 := uint8(num & 0xFF000000 >> 24)
	byteSlice := []byte{index0, index8, index16, index24}
	return byteSlice
}
