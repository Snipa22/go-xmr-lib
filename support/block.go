package support

import (
	"encoding/binary"
	"encoding/hex"
	"github.com/Snipa22/go-xmr-lib/support/crypto"
	"github.com/Snipa22/go-xmr-lib/support/serialization"
)

// Things to be implemented
// 1. convert_blob -> Converts a block template blob that's been modified into a hashable object
// 1.1 parse_and_validate_block_from_blob -> Created a Block struct from a blob of data provided by the Get Block Template RPC call, that is then modified to suit usages
// 1.2 get_block_hashing_blob -> Converts the blob into a block hashing blob

/*
ParseBlockFromTemplateBlob

Given a block template from the daemon, this will parse it into a usable Block object for handling data properly in go structs
*/
func ParseBlockFromTemplateBlob(blob string) (serialization.Block, error) {
	var b serialization.Block
	blobInBytes, err := hex.DecodeString(blob)
	if err != nil {
		return b, err
	}
	// Get the Major Version, uint8
	val, blobInBytes, err := serialization.ReadUint(blobInBytes)
	if err != nil {
		return b, err
	}
	b.MajorVersion = uint8(val)

	// Get the Minor Version, uint8
	val, blobInBytes, err = serialization.ReadUint(blobInBytes)
	if err != nil {
		return b, err
	}
	b.MinorVersion = uint8(val)

	// Get the Timestamp, uint64
	val, blobInBytes, err = serialization.ReadUint(blobInBytes)
	if err != nil {
		return b, err
	}
	b.Timestamp = val

	// Get the previous hash, which is an array of 32 bytes in uint8 form, stored as 32 bytes in the array
	bytesCopied := copy(b.PreviousID[:], blobInBytes[0:32])
	blobInBytes = blobInBytes[bytesCopied:]

	// Get the nonce, uint32, but is stored as a block of 4 bytes...  Jackassery.
	b.Nonce = binary.BigEndian.Uint32(blobInBytes[0:4])
	blobInBytes = blobInBytes[4:]

	// Start Transaction Processing (Miner Transaction)
	var t serialization.Transaction

	// Get Version, uint64
	val, blobInBytes, err = serialization.ReadUint(blobInBytes)
	if err != nil {
		return b, err
	}
	t.Version = val

	// Get UnlockTime, uint64 -- Could be a timestamp OR a block ID
	val, blobInBytes, err = serialization.ReadUint(blobInBytes)
	if err != nil {
		return b, err
	}
	t.UnlockTime = val

	// Start processing the t.vin fields
	// These are the Variant In fields.

	// Move forwards by 1 as the array is one object in length
	blobInBytes = blobInBytes[1:]

	// Move forwards by 1 as the key type for the array is going to be the TransactionInGenesis type
	blobInBytes = blobInBytes[1:]

	// Load the resulting blob data into the correct portion of Transaction/TransactionsIn
	var tig serialization.TransactionInGenesis
	var ti serialization.TransactionIn

	// Get the genesis height
	val, blobInBytes, err = serialization.ReadUint(blobInBytes)
	if err != nil {
		return b, err
	}
	tig.Height = val
	tig.Used = true
	ti.Genesis = tig

	t.TransactionsIn = append(t.TransactionsIn, ti)

	// Move forwards by 1 as the array is one object in length
	blobInBytes = blobInBytes[1:]

	// Load the resulting blob data into the correct portion of Transaction/TransactionsOut
	var to serialization.TransactionOut

	val, blobInBytes, err = serialization.ReadUint(blobInBytes)
	if err != nil {
		return b, err
	}
	to.Amount = val

	// Outbound type is to key, or 0x02.  Increment bytes by one
	blobInBytes = blobInBytes[1:]

	// Write a new TransactionOutToKey
	var totk serialization.TransactionOutToKey
	bytesCopied = copy(totk.PublicKey[:], blobInBytes[0:32])
	blobInBytes = blobInBytes[bytesCopied:]
	totk.Used = true
	to.Key = totk

	t.TransactionsOut = append(t.TransactionsOut, to)

	// Get the number of bytes to read into "extra"
	val, blobInBytes, err = serialization.ReadUint(blobInBytes)
	if err != nil {
		return b, err
	}

	// With val set to the # of bytes to read, slice and go
	t.Extra = blobInBytes[0:val]
	blobInBytes = blobInBytes[val:]

	// RingCT Type is 0.  Advance one byte
	blobInBytes = blobInBytes[1:]

	// Miner Transaction Complete!  Load to the main storage.  AYE AYE CAPTAIN!
	b.MinerTxn = t

	// Get the number of hashes in the tx_hashes field
	val, blobInBytes, err = serialization.ReadUint(blobInBytes)
	if err != nil {
		return b, err
	}

	// Attempt to get <val> hashes and append to the main store
	for ; val > 0; val-- {
		var iSlice [32]byte
		bytesCopied = copy(iSlice[:], blobInBytes[0:32])
		b.TxnHashes = append(b.TxnHashes, iSlice)
		blobInBytes = blobInBytes[bytesCopied:]
	}

	return b, nil
}

/*
GetBlockHashingBlob

Given a block object from the serialization library, it will convert the header of the object into a miner-hashable object
*/
func GetBlockHashingBlob(b serialization.Block) ([]byte, error) {
	// Source: cryptonote_format_utils.cpp
	// Original: get_block_hashing_blob(const block& b, blobdata& blob) - Line: 678-ish
	var sbh []byte = b.BlockHeader.Serialize()
	bmt := getBlockMerkleTreeHash(b)
	sbh = append(sbh, bmt[:]...)
	sbh = serialization.WriteUint(sbh, uint64(len(b.TxnHashes)+1))
	return sbh, nil
}

/*
GetBlockID

Given a block object, hash the blob to pull the block ID for the object, may be deprecated by the implementation of a proper response from the daemon
*/
func GetBlockID(b serialization.Block) [32]byte {
	return crypto.KeccakOneShot(b.GetBlob())
}

// Supporting hash systems
/*
getBlockMerkleTreeHash

Given a block object, this converts the transaction tree into a merkle hash tree for usage in the hashing blob
*/
func getBlockMerkleTreeHash(b serialization.Block) [32]byte {
	// Original: get_tx_tree_hash(const block& b) - cryptonote_format_utils.cpp:875-ish

	// Get the transaction hash.
	var txList [][32]byte
	txList = append(txList, getTransactionHash(b.MinerTxn))

	// Shift the hashes into a new slice, first one is the txn hash, then add all other hashes to the end.
	txList = append(txList, b.TxnHashes...)

	// Get the tree hash, this is the return.  Need to abstract some of this to a support library...
	return crypto.TreeHash(txList)
}
