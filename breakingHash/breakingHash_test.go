package breakingHash

import (
	"math/rand"
	"testing"
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
	result := BruteForceLookup2(testHash, 1)
	if result != testString {
		t.Errorf("BruteForceLookup2(%d, 1) = %v; want %v", testHash, result, testString)
	}
}
