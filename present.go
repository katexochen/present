// Package present implements the ultra-lightweight block cipher PRESENT as defined by Bogdanov et al. [1].
//
// 1. Bogdanov A. et al. (2007) PRESENT: An Ultra-Lightweight Block Cipher.
// In: Paillier P., Verbauwhede I. (eds) Cryptographic Hardware and Embedded Systems - CHES 2007.
// CHES 2007. Lecture Notes in Computer Science, vol 4727. Springer, Berlin, Heidelberg
package present

import (
	"strconv"
)

// keySizeError represents an invalid PRESENT key length.
type keySizeError int

func (k keySizeError) Error() string {
	return "present: invalid key size " + strconv.Itoa(int(k))
}

// NewCipher creates a new Block.
// The argument should be the PRESENT key,
// which is either 10 or 16 bytes long
// for key lengths of 80 bits and 128 bits respectively.
func NewCipher(key []byte, rounds int) (*Block, error) {
	switch len(key) {
	case 10:
		b := &Block{
			keySize:   len(key),
			roundKeys: make([]uint64, rounds+1),
			numRounds: rounds,
		}
		return b, b.SetKey(key)
	case 16:
		b := &Block{
			keySize:   len(key),
			roundKeys: make([]uint64, rounds+1),
			numRounds: rounds,
		}
		return b, b.SetKey(key)
	default:
		return nil, keySizeError(len(key))
	}
}
