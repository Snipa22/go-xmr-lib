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
	PublicKey [32]byte
	Used      bool
}

func (totk TransactionOutToKey) Serialize() []byte {
	var s []byte = []byte{0x02}
	s = append(s, totk.PublicKey[:]...)
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
	Extra           []uint8
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
	s = WriteUint(s, uint64(len(tp.Extra)))
	s = append(s, tp.Extra...)
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
