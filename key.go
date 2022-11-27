package present

// key represents an implementation of the PRESENT key schedule for deriving the round keys.
type key interface {
	rotate()
	sBox()
	xor(ctr uint64)
	roundKey() uint64
}

// decompose converts an 80-bit or 128-bit key into a pair of 64-bit integers.
func decompose(key []byte) (A, B uint64) {
	for i, x := range key {
		if i < 8 {
			shift := 56 - i*8
			A |= uint64(x) << uint64(shift)
		} else {
			shift := 120 - i*8
			B |= uint64(x) << uint64(shift)
		}
	}
	return
}

// updateKey updates key based on the key schedule and current round counter.
func updateKey(k key, ctr int) {
	k.rotate()
	k.sBox()
	k.xor(uint64(ctr))
}

// expandKey calculates the round keys and writes them into the given slice.
func expandKey(k key, numRounds int, dstRoundKeys []uint64) {
	for ctr := 0; ctr < numRounds; ctr++ {
		dstRoundKeys[ctr] = k.roundKey()
		updateKey(k, ctr+1)
	}
	dstRoundKeys[numRounds] = k.roundKey()
}
