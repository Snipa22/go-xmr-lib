package support

import (
	"github.com/Snipa22/go-xmr-lib/support/crypto"
	"github.com/Snipa22/go-xmr-lib/support/serialization"
)

func getTransactionPrefixHash(t serialization.Transaction) [32]byte {
	// Given a Transaction t, extract the TransactionPrefix TP and serialize it.
	// Given the resulting serialized data, cn_fast_hash (keccak-256) it.
	hash := crypto.KeccakOneShot(t.TransactionPrefix.Serialize())
	return hash
}

func getTransactionHash(t serialization.Transaction) [32]byte {
	// Original source: cryptonote_format_utils.cpp:617-ish
	// With hashes be three, may thee get the result thoust desire.
	var hs [3][32]byte

	// Thou must take tine prefix, and hash it!
	// Original : get_transaction_prefix_hash(t (Transaction), hashes[0] (crypto::hash))
	hs[0] = getTransactionPrefixHash(t)

	// Base RingCT Transaction Hash Data - byte 0 for the main txn, due to RingCTType being null (0x0)
	// So we're gonna short-cut this...
	// hexHash, _ := hex.DecodeString("044852B2A670ADE5407E78FB2863C51DE9FCB96542A07186FE3AEDA6BB8A116D")
	// copy(hs[1][:], hexHash[0:32])
	hs[1] = crypto.KeccakOneShot([]byte{0})

	// Null hashes are value 0, with no additional data
	hs[2] = [32]byte{0}

	var ah []byte
	ah = append(ah, hs[0][:]...)
	ah = append(ah, hs[1][:]...)
	ah = append(ah, hs[2][:]...)
	var h [32]byte = crypto.KeccakOneShot(ah)
	return h
}
