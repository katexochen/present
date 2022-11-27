package present

import (
	"fmt"
	"os"
	"testing"
)

func Test_sBoxLayer(t *testing.T) {
	t.Run("s box layer", func(t *testing.T) {
		var state uint64 = 0x123456789abcdef
		var expected uint64 = 0xc56b90ad3ef84712
		actual := SBoxLayer(state, SBox)
		if actual != expected {
			t.Fail()
		}
	})
	t.Run("inverse s box layer", func(t *testing.T) {
		var state uint64 = 0xc56b90ad3ef84712
		var expected uint64 = 0x123456789abcdef
		actual := SBoxLayer(state, SBoxInv)
		if actual != expected {
			t.Fail()
		}
	})
}

func Test_pLayer(t *testing.T) {
	t.Run("p layer", func(t *testing.T) {
		var state uint64 = 0xaaaaaaaaaaaaaaaa
		var expected uint64 = 0xffff0000ffff0000
		actual := PLayer(state, P)
		if actual != expected {
			fmt.Fprintf(os.Stderr, "expected: %x\ngot     : %x\n", expected, actual)
			t.Fail()
		}
	})
	t.Run("inverse p layer", func(t *testing.T) {
		var state uint64 = 0xffff0000ffff0000
		var expected uint64 = 0xaaaaaaaaaaaaaaaa
		actual := PLayer(state, PInv)
		if actual != expected {
			fmt.Fprintf(os.Stderr, "expected: %x\ngot     : %x\n", expected, actual)
			t.Fail()
		}
	})
}
