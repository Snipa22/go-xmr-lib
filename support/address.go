package support

import (
	"bytes"
	"github.com/Snipa22/go-xmr-lib/support/Base58"
)

func ValidateAddress(address string, tag []byte) (bool, error) {
	addr, err := Base58.DecodeFromString(address)
	if err != nil {
		return false, err
	}
	if bytes.Compare(addr[0:len(tag)], tag) != 0 {
		return false, nil
	}
	return true, nil
}
