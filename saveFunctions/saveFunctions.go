package saveFunctions

import (
	"fmt"
	"math"
	"strings"
	"unicode"

	"github.com/abramtrinh/wc3-save-code/lookup2"
)

/*
	https://jass.sourceforge.net/doc/types.shtml
	JASS Types:
	integer - int32
	real - 32-bit IEEE-754 AKA float32
	arrays - a sparse hashtable with a fixed size of 8192
*/

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

func CaseHash(key string) int32 {
	//TODO: error handling, might not need? just test, handle valid key? might just make function

	result := 0
	for _, runeValue := range key {
		result *= 2
		if unicode.IsUpper(runeValue) {
			result += 1
		}
	}
	return int32(result)
}

func Key2ParityKey(key string) int32 {
	// TODO: error handling here too for func and handle error from stringhash, test valid key
	sHash, _ := StringHash(key)
	cHash := CaseHash(key)
	var result int32 = sHash + cHash

	if result < 0 {
		result *= -1
	}

	return (result % 100)

}

////////////////////////////////////////////////////////////////////////////////////////////////////
// Everything below here is very much spooky and weird. Exported globals and stuff.
////////////////////////////////////////////////////////////////////////////////////////////////////

var Savecode_F int32 = 0
var Savecode_I int32 = 0
var Savecode_V map[int32]int32 = make(map[int32]int32)

var Savecode_Digits map[int32]float32 = make(map[int32]float32)
var Savecode_BigNum map[int32]int32 = make(map[int32]int32)

var BigNum_F int32 = 0
var BigNum_I int32 = 0
var BigNum_V map[int32]int32 = make(map[int32]int32)

var BigNum_List map[int32]int32 = make(map[int32]int32)
var BigNum_Base map[int32]int32 = make(map[int32]int32)

var BigNum_LF int32 = 0
var BigNum_LI int32 = 0
var BigNum_LV map[int32]int32 = make(map[int32]int32)
var BigNum_LLeaf map[int32]int32 = make(map[int32]int32)
var BigNum_LNext map[int32]int32 = make(map[int32]int32)

// function s__Savecode__allocate takes nothing returns integer
// local integer sc=s__Savecode__allocate()
func SaveCodeAllocate() int32 {
	temp := Savecode_F
	if temp != 0 {
		Savecode_F = Savecode_V[temp]
	} else {
		Savecode_I += 1
		temp = Savecode_I
	}
	if temp > 8190 {
		return 0
	}
	Savecode_V[temp] = -1
	return temp
}

// function s__BigNum__allocate takes nothing returns integer
// local integer b=s__BigNum__allocate()
func BigNumAllocate() int32 {
	temp := BigNum_F
	if temp != 0 {
		BigNum_F = BigNum_V[temp]
	} else {
		BigNum_I += 1
		temp = BigNum_I
	}
	if temp > 8190 {
		return 0
	}
	BigNum_V[temp] = -1
	return temp
}

// function s__BigNum_create takes integer base returns integer
// set s__Savecode_bignum[sc]=s__BigNum_create(VJSELogic___BASE())
func BigNumCreate(base int32) int32 {
	b := BigNumAllocate()
	BigNum_List[b] = 0
	BigNum_Base[b] = base
	return b
}

/*
constant string VJSELogic___CHARSET="0123456789@abcdefghijkmnopqrstuvwxyz#ABCDEFGHJKLMNOPQRSTUVWXYZ%"

VJSELogic___BASE()					==
VJSELogic___charsetlen()			==
StringLength(VJSELogic___CHARSET)	==
63									==
10+25+25+3	(0-9, a-z(no l), A-Z(no I), @#%)
*/
func Base() int32 {
	return 63
}

// function s__Savecode_create takes nothing returns integer
// local integer theCode=s__Savecode_create()
func SaveCodeCreate() int32 {
	sc := SaveCodeAllocate()
	Savecode_Digits[sc] = 0
	Savecode_BigNum[sc] = BigNumCreate(Base())
	return sc
}

// function s__BigNum_l__allocate takes nothing returns integer
// local integer bl=s__BigNum_l__allocate()
func BigNumLAllocate() int32 {
	temp := BigNum_LF
	if temp != 0 {
		BigNum_LF = BigNum_LV[temp]
	} else {
		BigNum_LI += 1
		temp = BigNum_LI
	}
	if temp > 8190 {
		return 0
	}
	BigNum_LV[temp] = -1
	return temp
}

// function s__BigNum_l_create takes nothing returns integer
func BigNumLCreate() int32 {
	bl := BigNumLAllocate()
	BigNum_LLeaf[bl] = 0
	BigNum_LNext[bl] = 0
	return bl
}

// function s__BigNum_MulSmall takes integer this,integer x returns nothing
// call s__BigNum_MulSmall(s__Savecode_bignum[this],max+1)
func BigNumMulSmall(temp, x int32) {
	cur := BigNum_List[temp]
	var product, remainder, carry int32 = 0, 0, 0

	for !(cur == 0 && carry == 0) {
		product = x*BigNum_LLeaf[cur] + carry
		carry = product / BigNum_Base[temp]
		remainder = product - carry*BigNum_Base[temp]
		BigNum_LLeaf[cur] = remainder

		if BigNum_LNext[cur] == 0 && carry != 0 {
			BigNum_LNext[cur] = BigNumLCreate()
		}

		cur = BigNum_LNext[cur]
	}

}

// function s__BigNum_AddSmall takes integer this,integer carry returns nothing
// call s__BigNum_AddSmall(s__Savecode_bignum[this],val)
func BigNumAddSmall(temp, carry int32) {
	//next not used?
	//var next, cur, sum int32 = 0, BigNum_List[temp], 0
	var cur, sum int32 = BigNum_List[temp], 0

	if cur == 0 {
		cur = BigNumLCreate()
		BigNum_List[temp] = cur
	}

	for !(carry == 0) {
		sum = BigNum_LLeaf[cur] + carry
		carry = sum / BigNum_Base[temp]
		sum = sum - carry*BigNum_Base[temp]
		BigNum_LLeaf[cur] = sum

		if BigNum_LNext[cur] == 0 {
			BigNum_LNext[cur] = BigNumLCreate()
		}
		cur = BigNum_LNext[cur]
	}
}

// function s__Savecode_Encode takes integer this,integer val,integer max returns nothing
// call s__Savecode_Encode(myCode,parity,parity+1)
func SaveCodeEncode(temp, val, max int32) {
	//Log base Base() of (max+1)
	customLog := math.Log(float64(max+1)) / math.Log(float64(Base()))
	Savecode_Digits[temp] = Savecode_Digits[temp] + float32(customLog)

	BigNumMulSmall(Savecode_BigNum[temp], max+1)
	BigNumAddSmall(Savecode_BigNum[temp], val)
}

// function s__VJSE_Save takes player p,string playerKey,integer codeVersion returns integer
func Save(playerKey string, codeVersion int32) {
	myCode := SaveCodeCreate()
	parity := Key2ParityKey(playerKey)

	SaveCodeEncode(myCode, parity, parity+1)

}
