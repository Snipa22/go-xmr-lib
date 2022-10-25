package support

import (
	"bytes"
	"errors"
	"github.com/Snipa22/go-xmr-lib/support/base58"
	"golang.org/x/crypto/sha3"
)

var ErrInvalidChecksum = errors.New("invalid checksum for address")

func IsValidTestnet(address string) (valid bool, err error) {
	// Testnet tags - 35/3f
	valid, err = ValidateAddress(address, []byte{0x35})
	if valid || err != nil {
		return
	}
	valid, err = ValidateAddress(address, []byte{0x3f})
	if valid || err != nil {
		return
	}
	return false, nil
}

func IsValidMainnet(address string) (valid bool, err error) {
	// Testnet tags - 12/2a/13/11
	valid, err = ValidateAddress(address, []byte{0x2a})
	if valid || err != nil {
		return
	}
	valid, err = ValidateAddress(address, []byte{0x12})
	if valid || err != nil {
		return
	}
	valid, err = ValidateAddress(address, []byte{0x13})
	if valid || err != nil {
		return
	}
	valid, err = ValidateAddress(address, []byte{0x11})
	if valid || err != nil {
		return
	}
	return false, nil
}

func ValidateAddress(address string, tag []byte) (bool, error) {
	addr, err := base58.DecodeFromString(address)
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
