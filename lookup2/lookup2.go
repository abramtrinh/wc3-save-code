package lookup2

/*
Original C source referenced:
"lookup2.c" / "Newhash" / "My Hash"
http://burtleburtle.net/bob/hash/evahash.html
http://burtleburtle.net/bob/c/lookup2.c
http://burtleburtle.net/bob/hash/doobs.html

For more reading:
https://en.wikipedia.org/wiki/Jenkins_hash_function
*/

// Reversibly mixes the values of (a, b, c).
func mix(a, b, c uint32) (uint32, uint32, uint32) {
	// Values used to left & right bitshift
	shiftValues := [9]uint32{13, 8, 13, 12, 16, 5, 3, 10, 15}

	for i := 0; i < len(shiftValues); i += 3 {
		a -= b
		a -= c
		// Shifts 13, 12, 3
		a ^= (c >> shiftValues[i])

		b -= c
		b -= a
		// Shifts 8, 16, 10
		b ^= (a << shiftValues[i+1])

		c -= a
		c -= b
		// Shifts 13, 5, 15
		c ^= (b >> shiftValues[i+2])
	}

	return a, b, c

}

// Hashes a variable-length string key into an unsigned 32-bit value.
func Hash(key string, seed int) uint32 {
	length := len(key)

	// In Go, indexing strings yields bytes, not characters/runes.
	s := string2RuneSlice(key)

	var a, b, c uint32 = 0x9e3779b9, 0x9e3779b9, uint32(seed)

	var lenLeft, i int = length, 0

	for lenLeft >= 12 {
		a += s[i] + (s[i+1] << 8) + (s[i+2] << 16) + (s[i+3] << 24)
		b += s[i+4] + (s[i+5] << 8) + (s[i+6] << 16) + (s[i+7] << 24)
		c += s[i+8] + (s[i+9] << 8) + (s[i+10] << 16) + (s[i+11] << 24)
		a, b, c = mix(a, b, c)
		i += 12
		lenLeft -= 12
	}

	c += uint32(length)

	switch lenLeft {
	case 11:
		c += s[i+10] << 24
		fallthrough
	case 10:
		c += s[i+9] << 16
		fallthrough
	case 9:
		c += s[i+8] << 8
		fallthrough
	case 8:
		b += s[i+7] << 24
		fallthrough
	case 7:
		b += s[i+6] << 16
		fallthrough
	case 6:
		b += s[i+5] << 8
		fallthrough
	case 5:
		b += s[i+4]
		fallthrough
	case 4:
		a += s[i+3] << 24
		fallthrough
	case 3:
		a += s[i+2] << 16
		fallthrough
	case 2:
		a += s[i+1] << 8
		fallthrough
	case 1:
		a += s[i]
	}

	//Only returning c, so no need to assign a & b
	_, _, c = mix(a, b, c)

	return c
}

//Converts string to rune slice so I can use a "for loop" instead of a "for range loop" to iterate over it.
func string2RuneSlice(input string) []uint32 {
	var runeSlice []uint32

	// Using a "for loop" to iterate over strings in Go iterates the bytes not the runes.
	// This is an issue when you have runes >1 bytes, like Chinese characters.
	for _, runeValue := range input {
		runeSlice = append(runeSlice, uint32(runeValue))
	}
	return runeSlice
}
