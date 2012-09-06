package ezhik

import (
	"math/rand"
)

func Encode(data []byte, n int, seed int64) (block []byte) {
	blockLen := (len(data) + n - 1) / n

	block = make([]byte, blockLen)
	r := rand.New(rand.NewSource(seed))
	for i := 0; i < n; i++ {
		if r.Intn(2) == 1 {
			var cur []byte
			if i < n-1 {
				cur = data[i*blockLen : (i+1)*blockLen]
			} else {
				cur = make([]byte, blockLen)
				// In case of small files, it's possible that last few blocks have zero length
				// For example, len(data) = 425, n = 128, blockLen = 4, 
				// but 4 * 128 = 512.
				if i*blockLen < len(data) {
					copy(cur, data[i*blockLen:])
				}
			}
			for i, v := range cur {
				block[i] ^= v
			}
		}
	}
	return
}
