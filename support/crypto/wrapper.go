package crypto

import "golang.org/x/crypto/sha3"

func KeccakOneShot(in []byte) [32]byte {
	keccak := sha3.NewLegacyKeccak256()
	keccak.Write(in)
	var data [32]byte
	copy(data[:], keccak.Sum([]byte{})[0:32])
	return data
}
