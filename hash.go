package ttlmap

type hashValue uint64

// FNV hash based on gopacket.
// See http://isthe.com/chongo/tech/comp/fnv/.
func fnvHash(s []byte) (h hashValue) {
	h = fnvBasis
	for i := 0; i < len(s); i++ {
		h ^= hashValue((s)[i])
		h *= fnvPrime
	}
	return
}

const fnvBasis = 14695981039346656037
const fnvPrime = 1099511628211
