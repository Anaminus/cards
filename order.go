package cards

import (
	crand "crypto/rand"
	"math/big"
)

// PerfectShuffle shuffles a group of cards with perfect randomness by using
// the Fisher-Yates shuffle.
func PerfectShuffle(g Group) {
	for i := g.Len() - 1; i > 0; i-- {
		j, _ := crand.Int(crand.Reader, big.NewInt(int64(i)+1))
		g.Swap(i, int(j.Int64()))
	}
}

// Reverse reverses the cards in a group.
func Reverse(g Group) {
	for i := g.Len()/2 - 1; i >= 0; i-- {
		j := g.Len() - 1 - i
		g.Swap(i, j)
	}
}
