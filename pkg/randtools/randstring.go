package randtools

import (
	cryptorand "crypto/rand"
	"math/rand/v2"
	"sync"
	"unsafe"
)

const (
	Alphabet    = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	AlphabetLen = 62

	letterIndexBits = 6
	letterIndexMask = 1<<letterIndexBits - 1
	letterIndexMax  = 63 / letterIndexMask
)

var (
	mu     sync.Mutex
	source *rand.Rand
)

func init() {
	var seed [32]byte
	cryptorand.Read(seed[:])
	chaCha8 := rand.NewChaCha8(seed)
	source = rand.New(chaCha8)
}

func DefaultGenerateString(stringLen int) string {
	return GenerateString(source, Alphabet, AlphabetLen, stringLen)
}

// This function is thread-safe
func GenerateString(source rand.Source, alphabet string, alphabetLen, stringLen int) string {
	result := make([]byte, stringLen)

	mu.Lock()
	cache := source.Uint64()
	mu.Unlock()

	remain := letterIndexMax

	i := stringLen - 1
	for i >= 0 {
		if remain == 0 {
			mu.Lock()
			cache = source.Uint64()
			mu.Unlock()

			remain = letterIndexMax
		}

		idx := int(cache & letterIndexMask)
		if idx < alphabetLen {
			result[i] = alphabet[idx]
			i--
		}

		cache >>= letterIndexBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&result))
}
