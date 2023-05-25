package saveFunctions

import (
	"math"
	"strings"
)

////////////////////////////////////////////////////////////////////////////////////////////////////
// Everything below here is very much spooky and weird. Exported globals and stuff.
// Please forgive me for what I have done.
// Implementation is as close to JASS's script file as possible, not based on idiomatic Go.
// Should be refactored after confirming it works exactly like in JASS.
// Currently blocked by JASS's SetRandomSeed & GetRandomInt since I don't have source code to work off of.
// Function that has this issue is SaveCodeObfuscate(). Good luck.
////////////////////////////////////////////////////////////////////////////////////////////////////

const CONFIG_MAX_VERSION_SUPPORT int32 = 150

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
	var product, remainder, carry int32 = math.MinInt32, math.MinInt32, 0

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
	var cur, sum int32 = BigNum_List[temp], math.MinInt32

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

func HashN() int32 {
	return 5000
}

func BigNumLDeallocate(temp int32) {
	if temp == 0 {
		return
	} else if BigNum_LV[temp] != -1 {
		return
	}
	//Below not implemented cause function does nothing
	//call s__BigNum_l_onDestroy(this)
	BigNum_LV[temp] = BigNum_LF
	BigNum_LF = temp
}

func BigNumLClean(temp int32) bool {
	if BigNum_LNext[temp] == 0 && BigNum_LLeaf[temp] == 0 {
		return true
	} else if BigNum_LNext[temp] != 0 && BigNumLClean(BigNum_LNext[temp]) {
		BigNumLDeallocate(BigNum_LNext[temp])
		BigNum_LNext[temp] = 0
		return BigNum_LLeaf[temp] == 0
	} else {
		return false
	}
}

func BigNumClean(temp int32) {
	cur := BigNum_List[temp]
	BigNumLClean(cur)
}

func SaveCodeClean(temp int32) {
	BigNumClean(Savecode_BigNum[temp])
}

func SaveCodeHash(temp int32) int32 {
	var lHash, x, cur int32 = 0, 0, BigNum_List[Savecode_BigNum[temp]]

	if !(cur == 0) {
		x = BigNum_LLeaf[cur]
		//function ModuloInteger takes integer dividend, integer divisor returns integer
		var dividend int32 = lHash + 79*lHash/(x+1) + 293*x/(1+lHash-(lHash/Base()))*Base() + 479
		//VJSELogic___HASHN = 5000
		var divisor int32 = HashN()
		lHash = dividend % divisor
		cur = BigNum_LNext[cur]
	}

	return lHash
}

func SaveCodeLength(temp int32) float32 {
	return Savecode_Digits[temp]
}

func SaveCodePad(temp int32) {
	cur := BigNum_List[Savecode_BigNum[temp]]
	var prev int32
	var maxlen int32 = int32(1.0 + SaveCodeLength(temp))

	for cur != 0 {
		prev = cur
		cur = BigNum_LNext[cur]
		maxlen--
	}

	for maxlen > 0 {
		BigNum_LNext[prev] = BigNumLCreate()
		prev = BigNum_LNext[prev]
		maxlen--
	}
}

func CharToI(c rune) int32 {
	charSet := "0123456789@abcdefghijkmnopqrstuvwxyz#ABCDEFGHJKLMNOPQRSTUVWXYZ%"

	index := strings.IndexRune(charSet, c)

	return int32(index)
}

// function VJSELogic___scommhash takes string s returns integer
// local integer key=VJSELogic___scommhash(p)+loadtype*73
func SCommHash(s string) int32 {
	var count map[int32]int32 = make(map[int32]int32)
	var x int32 = 0

	s = strings.ToUpper(s)

	for _, runeValue := range s {
		x = CharToI(runeValue)
		count[x] = count[x] + 1
	}
	var i, len int32 = 0, Base()
	x = 0
	for ; i < len; i++ {
		x = count[i]*count[i]*i + count[i]*x + x + 199
	}

	if x < 0 {
		x *= -1
	}
	return x
}

func SaveCodeObfuscate(temp, key, sign int32) {
	// BLOCKED HERE; CONTINUE IF YOU FIND A SOLUTION TO GetRandomInt and SetRandomSeed
}

// function s__Savecode_Save takes integer this,string p,integer loadtype returns string
// set s__VJSE___PROC_P_SAVECODE=s__Savecode_Save(myCode,playerKey,s__VJSE___CONFIG_MAP_KEY)
func SaveCodeSave(temp int32, p string, loadtype int32) string {
	// Uncomment later.
	//key := SCommHash(p) + loadtype*73
	var lHash int32

	SaveCodeClean(temp)
	lHash = SaveCodeHash(temp)
	SaveCodeEncode(temp, lHash, HashN())
	SaveCodeClean(temp)
	SaveCodePad(temp)

	// Functions that still need to be implemented.
	//SaveCodeObfuscate(temp,key,1)
	//return SaveCodeToString

	//The return of this is the actual save code that you can use to load.
	return ""
}

// function s__VJSE_Save takes player p,string playerKey,integer codeVersion returns integer
func Save(playerKey string, codeVersion int32) {
	myCode := SaveCodeCreate()
	parity := Key2ParityKey(playerKey)

	SaveCodeEncode(myCode, parity, parity+1)
	//Did not implement HaveSavedHandle function.
	//It deallocates some of the global variables which shouldn't be an issue (if I run new state each time)

	//Did not implement FireEvent
	//Checks value of global variables.

	SaveCodeEncode(myCode, codeVersion, CONFIG_MAX_VERSION_SUPPORT)

	// THE save code.
	//saveCode := SaveCodeSave

}
