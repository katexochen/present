package present

var (
	// SBox is the substitution box of PRESENT.
	SBox = []byte{
		0xc, 0x5, 0x6, 0xb, 0x9, 0x0, 0xa, 0xd,
		0x3, 0xe, 0xf, 0x8, 0x4, 0x7, 0x1, 0x2,
	}

	// SBoxInv is the inverse substitution box of PRESENT.
	SBoxInv = []byte{
		0x5, 0xe, 0xf, 0x8, 0xc, 0x1, 0x2, 0xd,
		0xb, 0x4, 0x6, 0x3, 0x0, 0x7, 0x9, 0xa,
	}

	// P is the permutation of PRESENT.
	P = []byte{
		0, 16, 32, 48, 1, 17, 33, 49,
		2, 18, 34, 50, 3, 19, 35, 51,
		4, 20, 36, 52, 5, 21, 37, 53,
		6, 22, 38, 54, 7, 23, 39, 55,
		8, 24, 40, 56, 9, 25, 41, 57,
		10, 26, 42, 58, 11, 27, 43, 59,
		12, 28, 44, 60, 13, 29, 45, 61,
		14, 30, 46, 62, 15, 31, 47, 63,
	}

	// PInv is the inverse permutation of PRESENT.
	PInv = []byte{
		0, 4, 8, 12, 16, 20, 24, 28,
		32, 36, 40, 44, 48, 52, 56, 60,
		1, 5, 9, 13, 17, 21, 25, 29,
		33, 37, 41, 45, 49, 53, 57, 61,
		2, 6, 10, 14, 18, 22, 26, 30,
		34, 38, 42, 46, 50, 54, 58, 62,
		3, 7, 11, 15, 19, 23, 27, 31,
		35, 39, 43, 47, 51, 55, 59, 63,
	}
)

type Block struct {
	roundKeys []uint64
	keySize   int
	numRounds int
}

func (b *Block) SetKey(key []byte) error {
	if len(key) != b.keySize {
		return keySizeError(len(key))
	}

	switch b.keySize {
	case 10:
		expandKey(newKey80(key), b.numRounds, b.roundKeys)
	case 16:
		expandKey(newKey128(key), b.numRounds, b.roundKeys)
	default:
		return keySizeError(len(key))
	}

	return nil
}

func (b *Block) Encrypt(m uint64) uint64 {
	for i := 0; i < b.numRounds; i++ {
		m ^= b.roundKeys[i]
		m = SBoxLayer(m, SBox)
		m = PLayer(m, P)
	}
	roundKey := b.roundKeys[b.numRounds]
	m ^= roundKey
	return m
}

func (b *Block) Decrypt(c uint64) uint64 {
	c ^= b.roundKeys[b.numRounds]
	for i := b.numRounds - 1; i >= 0; i-- {
		c = PLayer(c, PInv)
		c = SBoxLayer(c, SBoxInv)
		c ^= b.roundKeys[i]
	}
	return c
}

func SBoxLayer(state uint64, s []byte) (result uint64) {
	for i := 0; i < 16; i++ {
		shift := 4 * uint(i)
		var mask uint64 = 0xF << shift
		x := (state & mask) >> shift
		y := uint64(s[x])
		z := y << shift
		result |= z
	}
	return
}

func PLayer(state uint64, p []byte) (result uint64) {
	for i := 0; i < len(p); i++ {
		result |= (state & 0x1) << p[i]
		state = state >> 1
	}
	return
}
