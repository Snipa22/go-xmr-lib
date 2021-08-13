package crypto

import "golang.org/x/crypto/sha3"

func TreeHash(hs [][32]byte) [32]byte {
	hash := [32]byte{}
	count := uint64(len(hs))
	if len(hs) == 1 {
		return hs[0]
	} else if len(hs) == 2 {
		h := sha3.NewLegacyKeccak256() // Hash Object
		h.Write(hs[0][:])
		h.Write(hs[1][:])
		var tempHash []byte
		tempHash = h.Sum(tempHash)
		copy(hash[:], tempHash[0:32])
	} else {
		// Borrowed from the tree-hash.c inmplenm
		var i, j uint64
		var cnt uint64 = count - 1
		var ints [][32]byte
		for i = 1; i < 8<<3; i <<= 1 {
			cnt |= cnt >> i
		}
		cnt &= ^(cnt >> 1)
		ints = make([][32]byte, cnt)
		for c, e := range hs {
			if uint64(c) >= 2*cnt-count {
				break
			}
			ints[c] = e
		}
		for i, j = 2*cnt-count, 2*cnt-count; j < cnt; i, j = i+2, j+1 {
			// cn_fast_hash(hashes[i], 2 * HASH_SIZE, ints[j]);
			// Input, Length, Result
			// Resut = KeccakOneShot(Input)
			// These are doubles, not singles.
			intHash := append(hs[i][:], hs[i+1][:]...)
			ints[j] = KeccakOneShot(intHash)
		}
		for {
			if cnt <= 2 {
				break
			}
			cnt >>= 1
			for i, j = 0, 0; j < cnt; i, j = i+2, j+1 {
				intHash := append(ints[i][:], ints[i+1][:]...)
				ints[j] = KeccakOneShot(intHash)
			}
		}
		hash = KeccakOneShot(append(ints[0][:], ints[1][:]...))
	}
	return hash
}
