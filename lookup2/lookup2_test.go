package lookup2

import (
	"fmt"
	"testing"
)

func TestHash(t *testing.T) {
	var tests = []struct {
		key  string
		seed int
		hash uint32
	}{
		{"", 0, 3175731469},
		{" ", 0, 2658412151},
		{"~", 0, 2374639685},
		{"E", 0, 597637742},
		{"AA", 0, 4050291262},
		{"Eb", 0, 2867312368},
		{"This is the time for all good men to come to the aid of their country", 0, 3481751101},
		{"This is the time for all good men to come to the aid of their country", 11, 339338997},
		//Tests multi-byte runes/characters below
		{"Ç", 0, 3014993710},
		{"■", 0, 1142109739},
		{"日本語", 0, 2971310703},
		{"Hello, 世界", 0, 2497382974},
	}

	for _, test := range tests {
		testName := fmt.Sprintf("Key: %v Seed: %v Hash: %v", test.key, test.seed, test.hash)
		t.Run(testName, func(t *testing.T) {
			testResult := Hash(test.key, test.seed)
			if testResult != test.hash {
				t.Errorf("%v %v returned %v, expected %v ", test.key, test.seed, testResult, test.hash)
			}
		})
	}
}
