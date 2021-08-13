package serialization

import (
	"encoding/binary"
	"errors"
)

var (
	errorBadBlob     error = errors.New("bad blob provided, unable to continue")
	errorUintTooLong error = errors.New("bad blob provided, uint is unable to be decoded, unable to continue")
)

func ReadUint(b []byte) (uint64, []byte, error) {
	val, byteCount := binary.Uvarint(b)
	if byteCount == 0 {
		return 0, b, errorBadBlob
	}
	if byteCount <= 0 {
		return 0, b, errorUintTooLong
	}
	return val, b[byteCount:], nil
}

func WriteUint(b []byte, v uint64) []byte {
	var tempBlob []byte = make([]byte, 32)
	bytesWritten := binary.PutUvarint(tempBlob, v)
	b = append(b, tempBlob[0:bytesWritten]...)
	return b
}
