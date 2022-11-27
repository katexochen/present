package present

// BlockSize is the PRESENT block size in bytes.
const BlockSize = 8

// Substitution and permutation tables for PRESENT.
var (
	SBox    = []byte{0xC, 5, 6, 0xB, 9, 0, 0xA, 0xD, 3, 0xE, 0xF, 8, 4, 7, 1, 2}
	SBoxInv = []byte{5, 0xE, 0xF, 8, 0xC, 1, 2, 0xD, 0xB, 4, 6, 3, 0, 7, 9, 0xA}
	P       = []byte{
		0, 16, 32, 48, 1, 17, 33, 49, 2, 18, 34, 50, 3, 19, 35, 51, 4, 20, 36, 52, 5, 21, 37, 53, 6,
		22, 38, 54, 7, 23, 39, 55, 8, 24, 40, 56, 9, 25, 41, 57, 10, 26, 42, 58, 11, 27, 43, 59, 12,
		28, 44, 60, 13, 29, 45, 61, 14, 30, 46, 62, 15, 31, 47, 63,
	}
	PInv = []byte{
		0, 4, 8, 12, 16, 20, 24, 28, 32, 36, 40, 44, 48, 52, 56, 60, 1, 5, 9, 13, 17, 21, 25, 29, 33,
		37, 41, 45, 49, 53, 57, 61, 2, 6, 10, 14, 18, 22, 26, 30, 34, 38, 42, 46, 50, 54, 58, 62, 3, 7,
		11, 15, 19, 23, 27, 31, 35, 39, 43, 47, 51, 55, 59, 63,
	}
)

type Block struct {
	roundKeys []uint64
	keySize   int
	numRounds int
}

func (b *Block) BlockSize() int {
	return BlockSize
}

func (b *Block) SetKey(key []byte) error {
	if len(key) != b.keySize {
		return KeySizeError(len(key))
	}

	switch b.keySize {
	case 10:
		expandKey(newKey80(key), b.numRounds, b.roundKeys)
	case 16:
		expandKey(newKey128(key), b.numRounds, b.roundKeys)
	default:
		return KeySizeError(len(key))
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
