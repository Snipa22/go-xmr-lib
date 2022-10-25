package serialization

// Transaction Inputs
type TransactionInGenesis struct {
	Height uint64
	Used   bool
}

func (tig TransactionInGenesis) Serialize() []byte {
	var s []byte = []byte{0xff}
	s = WriteUint(s, tig.Height)
	return s
}

type TransactionInToScript struct {
	PreviousHash   [32]byte
	PreviousOutput uint64
	SignatureSet   []uint8
	Used           bool
}

type TransactionInToScriptHash struct {
	PreviousHash   [32]byte
	PreviousOutput uint64
	Script         TransactionOutToScript
	SignatureSet   []uint8
	Used           bool
}

type TransactionInToKey struct {
	Amount     uint64    // Amount of coin transferred
	KeyOffsets []uint64  // Key offsets are numeric ID's
	KeyImage   [32]uint8 // Key Image is a hex
	Used       bool
}

type TransactionIn struct {
	Genesis    TransactionInGenesis
	Script     TransactionInToScript
	ScriptHash TransactionInToScriptHash
	Key        TransactionInToKey
}

func (ti TransactionIn) Serialize() []byte {
	var s []byte
	if ti.Genesis.Used {
		s = append(s, ti.Genesis.Serialize()...)
	}
	return s
}

// Transaction Outputs
type TransactionOutToScript struct {
	PublicKey [32]byte
	Script    []uint8
	Used      bool
}

type TransactionOutToScriptHash struct {
	Hash [32]byte
	Used bool
}

type TransactionOutToKey struct {
	Version   [1]byte
	PublicKey [32]byte
	Used      bool
	TaggedKey [1]byte
}

func (totk TransactionOutToKey) Serialize() []byte {
	var s = []byte{totk.Version[0]}
	s = append(s, totk.PublicKey[:]...)
	if totk.Version[0] >= 3 {
		s = append(s, totk.TaggedKey[0])
	}
	return s
}

type TransactionOut struct {
	Amount     uint64
	Script     TransactionOutToScript
	ScriptHash TransactionOutToScriptHash
	Key        TransactionOutToKey
}

func (to TransactionOut) Serialize() []byte {
	var s []byte
	s = WriteUint(s, to.Amount)
	if to.Key.Used {
		s = append(s, to.Key.Serialize()...)
	}
	return s
}

type TransactionPrefix struct {
	Version         uint64
	UnlockTime      uint64
	TransactionsIn  []TransactionIn
	TransactionsOut []TransactionOut
	Extra           TxExtra
}

func (tp TransactionPrefix) Serialize() []byte {
	// varint - Version
	// varint - unlocked_time
	// Vector store of the vin
	// Vector store of the vout
	// Slap on that extra data.  Mmmmm.  Extra.  Data.  Nomnom.
	var s []byte
	s = WriteUint(s, tp.Version)
	s = WriteUint(s, tp.UnlockTime)
	s = WriteUint(s, uint64(len(tp.TransactionsIn)))
	for _, e := range tp.TransactionsIn {
		s = append(s, e.Serialize()...)
	}
	s = WriteUint(s, uint64(len(tp.TransactionsOut)))
	for _, e := range tp.TransactionsOut {
		s = append(s, e.Serialize()...)
	}
	extraEnc := tp.Extra.Serialize()
	s = WriteUint(s, uint64(len(extraEnc)))
	s = append(s, extraEnc...)
	return s
}

type Transaction struct {
	TransactionPrefix
	Signatures [][32]byte
}

func (t Transaction) Serialize() []byte {
	// Cheating, as our ringCT sing type is 0
	var s []byte = t.TransactionPrefix.Serialize()
	s = append(s, 0)
	return s
}

type TxExtra struct {
	Padding        *[]byte // TX_EXTRA_TAG_PADDING - 0x00
	PubKey         []byte  // TX_EXTRA_TAG_PUBKEY - 0x01
	Nonce          []byte  // TX_EXTRA_NONCE - 0x02
	MergeMiningTag *[]byte // TX_EXTRA_MERGE_MINING_TAG - 0x03
	AddlPubKeys    *[]byte // TX_EXTRA_TAG_ADDITIONAL_PUBKEYS - 0x04
	MystMGTag      *[]byte // TX_EXTRA_MYSTERIOUS_MINERGATE_TAG - 0x0e
}

func (txe *TxExtra) Serialize() []byte {
	var s []byte
	if txe.Padding != nil {
		s = WriteUint(s, 0x00)
		s = append(s, *txe.Padding...)
	}
	s = WriteUint(s, 0x01)
	s = append(s, txe.PubKey...)
	s = WriteUint(s, 0x02)
	s = WriteUint(s, uint64(len(txe.Nonce)))
	s = append(s, txe.Nonce...)
	if txe.MergeMiningTag != nil {
		s = WriteUint(s, 0x03)
		s = append(s, *txe.MergeMiningTag...)
	}
	if txe.AddlPubKeys != nil {
		s = WriteUint(s, 0x04)
		s = append(s, *txe.AddlPubKeys...)
	}
	if txe.MystMGTag != nil {
		s = WriteUint(s, 0x0e)
		s = append(s, *txe.MystMGTag...)
	}
	return s
}

func (txe *TxExtra) UpdateNonce(in []byte) error {
	if len(in) == len(txe.Nonce) {
		txe.Nonce = in
	} else if len(in) > len(txe.Nonce) {
		return ErrNonceTooLong
	} else {
		newNonce := in
		for len(newNonce) < len(txe.Nonce) {
			newNonce = append(newNonce, 0)
		}
		txe.Nonce = newNonce
	}
	return nil
}

func ConstructTXExtra(in []byte) TxExtra {
	mutable := in
	extra := TxExtra{}
	for {
		// Perform length check at the start, not the end
		if len(mutable) == 0 {
			return extra
		}
		switch mutable[0] {
		case 0x00:
			extra.Padding = &[]byte{0}
			mutable = mutable[1:]
		case 0x01:
			extra.PubKey = append(extra.PubKey, mutable[1:33]...)
			mutable = mutable[33:]
		case 0x02:
			nonceLength := int(mutable[1])
			extra.Nonce = append(extra.Nonce, mutable[2:nonceLength+2]...)
			mutable = mutable[nonceLength+2:]
		}
	}
}
