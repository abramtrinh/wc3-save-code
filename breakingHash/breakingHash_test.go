package breakingHash

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/abramtrinh/wc3-save-code/lookup2"
)

// Keep 10 for iRf74ywBRf to be correct.
var stringLength = 10
var testString = "iRf74ywBRf"
var testHash uint32 = 771166165

func BenchmarkRandStringBytes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		randStringBytes(stringLength)
	}
}

func BenchmarkRandStringMask(b *testing.B) {
	for i := 0; i < b.N; i++ {
		randStringMask(stringLength)
	}
}

// If any of the global variables are changed. Test will take variable amount of time.
func TestBruteForceLookup2(t *testing.T) {
	// Deprecated
	rand.Seed(1)
	result, _ := BruteForceLookup2(testHash, 1)
	if result != testString {
		t.Errorf("BruteForceLookup2(%d, 1) = %v; want %v", testHash, result, testString)
	}
}

func TestUnHash(t *testing.T) {
	var tests = []struct {
		key string
	}{
		{""},
		{" "},
		{"~"},
		{"E"},
		{"AA"},
		{"Eb"},
		{"This is the time for all good men to come to the aid of their country"},
	}
	for _, test := range tests {
		testName := fmt.Sprintf("Key: %v", test.key)
		t.Run(testName, func(t *testing.T) {
			hash := lookup2.Hash(test.key, 0)
			unHashString := string(UnHash(hash))
			unHashReHash := lookup2.Hash(unHashString, 0)
			if hash != unHashReHash {
				t.Errorf("%v returned %v, expected %v ", test.key, unHashReHash, hash)
			}
		})
	}

}
