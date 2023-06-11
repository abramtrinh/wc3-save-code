package saveCode

type Save struct {
	name      string
	alphabet  string
	base      int
	max       int
	count     int
	char      map[int]string
	save      map[int]int
	temporary map[int]int
}

// Initializes the save code generator.
func NewSave(playerName string, maxValue int) *Save {
	fullAlphabet := "abcdefghkmnopqrstuvwxyzABCDEFGHJKLMNOPQRSTUVWXYZ0123456789"
	fullBase := len(fullAlphabet)

	p := new(Save)

	p.char = make(map[int]string)
	p.save = make(map[int]int)
	p.temporary = make(map[int]int)

	p.name = playerName
	p.max = maxValue
	p.count = 0

	for i := 1; i < p.max; i++ {
		//Indexing the string by the bytes since it is single-byte runepoints, not more.
		p.char[i] = string(fullAlphabet[i])
	}

	p.alphabet = string(fullAlphabet[0:1]) + string(fullAlphabet[p.max+1:fullBase])
	p.base = fullBase - p.max

	return p
}

// Sets up the internal state for the save code and returns the save code.
func (s *Save) SaveCode(rank int) string {
	// GetRandomInt(100, 500)
	s.temporary[1] = 499
	// GetRandomInt(501, 999)
	s.temporary[2] = 998
	s.save[s.count] = s.temporary[1]
	checksum := s.stringChecksum(s.name)
	s.temporary[0] = (s.temporary[2] - s.temporary[1]) * checksum
	s.count++
	s.save[s.count] = s.temporary[0]
	s.count++
	s.save[s.count] = s.temporary[0] * rank
	s.count++
	s.save[s.count] = s.temporary[2]

	return s.compile()
}

// Generates the save code string using the internal state.
func (s Save) compile() string {
	out := ""
	for i := 0; i <= s.count; i++ {
		x := s.encode(s.save[i])
		j := len(x)
		if j > 1 {
			out += s.char[j-1]
		}
		out += x
	}
	checksum := s.stringChecksum(s.name)
	out += s.encode(checksum)
	checksum = s.stringChecksum(out)
	out = s.encode(checksum) + out
	return out
}

// Encoding the str.
func (s Save) encode(i int) string {
	if i <= s.base {
		return string(s.alphabet[i])
	}
	str := ""
	b := 0
	for i > 0 {
		b = i - (i/s.base)*s.base
		str = string(s.alphabet[b]) + str
		i /= s.base
	}
	return str

}

// Converts the input string to a checksum.
func (s Save) stringChecksum(input string) int {
	checksum := 0
	for _, v := range input {
		temp := s.decode(string(v))
		checksum += temp
	}
	return checksum
}

// Decoding the str.
func (s Save) decode(str string) int {
	a := 0

	for len(str) != 1 {
		a = a*s.base + s.base*s.stringPosition(str[0:1])
		str = str[1:99]
	}

	return a + s.stringPosition(str)
}

// The index position of str
func (s Save) stringPosition(str string) int {
	for i, v := range s.alphabet {
		if str == string(v) {
			return i
		}
	}
	return -1
}
