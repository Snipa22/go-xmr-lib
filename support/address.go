package support

import (
	"bytes"
	"errors"
	"github.com/Snipa22/go-xmr-lib/support/Base58"
	"golang.org/x/crypto/sha3"
)

var ErrInvalidChecksum = errors.New("invalid checksum for address")

func ValidateAddress(address string, tag []byte) (bool, error) {
	addr, err := Base58.DecodeFromString(address)
	if err != nil {
		return false, err
	}
	keccak := sha3.NewLegacyKeccak256()
	keccak.Write(addr[0 : len(addr)-4])
	res := keccak.Sum([]byte{})
	if bytes.Compare(res[0:4], addr[len(addr)-4:]) != 0 {
		return false, ErrInvalidChecksum
	}
	if bytes.Compare(addr[0:len(tag)], tag) != 0 {
		return false, nil
	}
	return true, nil
}
