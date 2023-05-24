package saveFunctions

import (
	"fmt"
	"testing"
)

func TestStringHash(t *testing.T) {
	var tests = []struct {
		key  string
		hash int32
	}{
		{"a", -1587459251},
		{"A", -1587459251},
		{"hello", -1801350911},
		{"HELLO", -1801350911},
		{"hElLo", -1801350911},
		{"123", 1101291013},
		{"a1b2c3", -22837905},
	}

	for _, test := range tests {
		testName := fmt.Sprintf("Key: %v Hash: %v", test.key, test.hash)
		t.Run(testName, func(t *testing.T) {
			testResult, _ := StringHash(test.key)
			if testResult != test.hash {
				t.Errorf("%v returned %v, expected %v ", test.key, testResult, test.hash)
			}
		})
	}

	var testErrors = []struct {
		key string
	}{
		{" "},
		{""},
		{"Ç"},
		{"日本語"},
		{"#"},
	}

	for _, test := range testErrors {
		testName := fmt.Sprintf("Key: %v", test.key)
		t.Run(testName, func(t *testing.T) {
			_, err := StringHash(test.key)
			if err == nil {
				t.Errorf("%v should have errored. Did not error instead ", test.key)
			}
		})
	}
}
