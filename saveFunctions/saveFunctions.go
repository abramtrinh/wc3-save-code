package saveFunctions

import (
	"fmt"
	"strings"

	"github.com/abramtrinh/wc3-save-code/lookup2"
)

/*
	BattleTag Naming Policy
	https://us.battle.net/support/en/article/26963
	Need to support only alphanumeric (case-sensitive) characters.
	Supposedly supports accented characters.

	Tested nuances of StringHash() in Warcraft III: Reforged.
	All characters are string.ToUpper()'d if possible.
		e.g. ñ -> Ñ		б -> Б		a -> A
			StringHash() has the same hash value for both upper and lower case version.
	Forward slashes are converted to backslashes (but never used).
*/

// Takes in a valid alphanumeric string and hashes it into an int32 based on the WC3 implementation.
func StringHash(key string) (int32, error) {
	if len(key) <= 0 {
		return 0, fmt.Errorf("Error: nil input key.")
	}

	modKey := strings.ToUpper(key)

	// Only need to check for alphanumeric characters (uppercase only) since everything is ToUpper()'d.
	for _, runeValue := range modKey {
		// runeValue needs to have its values be between 48-57 || 65-90
		// Going with values outside of the range is does not StringHash correctly.
		if !(((runeValue >= 48) && (runeValue <= 57)) || ((runeValue >= 65) && (runeValue <= 90))) {
			// runeValue is not between 48-57 || 65-90
			return 0, fmt.Errorf("Error: %c is not a valid alphanumeric.", runeValue)
		}
	}

	// Cast the uint32 lookup2 hash result to int32 (which is what the JASS's StringHash does)
	return int32(lookup2.Hash(modKey, 0)), nil
}
