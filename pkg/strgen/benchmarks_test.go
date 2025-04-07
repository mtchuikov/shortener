package strgen

import (
	"math/rand/v2"
	"testing"
)

func BenchmarkSimpleRand(b *testing.B) {
	runeAlphabet := []rune(DefaultAlphabet)

	for range 10000 {
		b := make([]rune, 32)
		for i := range b {
			idx := rand.IntN(DefaultAlphabetLen)
			b[i] = runeAlphabet[idx]
		}

		_ = string(b)
	}
}

func BenchmarkGenerator(b *testing.B) {
	generage := New()

	for range 10000 {
		_ = generage.Generate(32)
	}
}
