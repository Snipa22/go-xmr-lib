package Base58

import (
	"encoding/binary"
	"errors"
	"strings"
)

const alphabet = `123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz`

var encodedBlockSizes = [9]int{0, 2, 3, 5, 6, 7, 9, 10, 11}

const fullBlockSize = 8
const encodedBlockSize = 11

var ErrInvalidBlockSize = errors.New("invalid block size for base58 encoding")
var ErrInvalidChar = errors.New("invalid character found in block decoding")
var ErrOverflow = errors.New("overflow found in block decoding")

/*
uint8BeTo64

This converts a set of uint8 objects (bytes) into an uint64
*/
func uint8BeTo64(in []byte) (uint64, error) {
	inLen := len(in)
	if inLen < 1 || inLen > fullBlockSize {
		return 0, ErrInvalidBlockSize
	}
	var retVal uint64 = 0
	caseVal := 9 - inLen
	offset := 0
	switch {
	case caseVal == 1:
		retVal |= uint64(in[offset])
		offset += 1
		fallthrough
	case caseVal == 2:
		retVal <<= 8
		retVal |= uint64(in[offset])
		offset += 1
		fallthrough
	case caseVal == 3:
		retVal <<= 8
		retVal |= uint64(in[offset])
		offset += 1
		fallthrough
	case caseVal == 4:
		retVal <<= 8
		retVal |= uint64(in[offset])
		offset += 1
		fallthrough
	case caseVal == 5:
		retVal <<= 8
		retVal |= uint64(in[offset])
		offset += 1
		fallthrough
	case caseVal == 6:
		retVal <<= 8
		retVal |= uint64(in[offset])
		offset += 1
		fallthrough
	case caseVal == 7:
		retVal <<= 8
		retVal |= uint64(in[offset])
		offset += 1
		fallthrough
	case caseVal == 8:
		retVal <<= 8
		retVal |= uint64(in[offset])
	}
	return retVal, nil
}

/*
encodeBlock

This performs the actual encoding of a block of data from the input encoders
*/
func encodeBlock(toEncode []byte) (string, error) {
	if len(toEncode) < 1 || len(toEncode) > fullBlockSize {
		return "", ErrInvalidBlockSize
	}
	res, err := uint8BeTo64(toEncode)
	if err != nil {
		return "", err
	}
	retVal := ""
	for 0 < res {
		retVal = string(alphabet[res%58]) + retVal
		res /= 58
	}
	return retVal, nil
}

/*
decodeBlock

This performs the actual decoding of a block object back to a byte slice
*/
func decodeBlock(toDecode string) ([]byte, error) {
	if len(toDecode) < 1 || len(toDecode) > encodedBlockSize {
		return nil, ErrInvalidBlockSize
	}
	blockSize := 0
	for k, v := range encodedBlockSizes {
		if v == len(toDecode) {
			blockSize = k
			break
		}
	}
	// Safety trap, though this should never be possible
	if blockSize <= 0 {
		return nil, ErrInvalidBlockSize
	}
	var resNum uint64 = 0
	var order uint64 = 1
	for i := range toDecode {
		i = len(toDecode) - 1 - i
		char := strings.Index(alphabet, string(toDecode[i]))
		if char == -1 {
			return nil, ErrInvalidChar
		}
		product := order*uint64(char) + resNum
		if product > 2^64 {
			return nil, ErrOverflow
		}
		resNum = product
		order *= 58
	}
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, resNum)
	return nil, nil
}

/*
EncodeToString

This function takes a byte array from the caller and encodes it in Monero Base58 format.
*/
func EncodeToString(src []byte) (string, error) {
	lenData := len(src)
	if lenData == 0 {
		return "", nil
	}
	fullBlockCount := lenData / fullBlockSize
	lastBlockSize := lenData % fullBlockSize
	retVal := ""
	for i := 0; i < fullBlockCount; i++ {
		res, err := encodeBlock(src[i*8 : (i+1)*8])
		if err != nil {
			return "", err
		}
		retVal += res
	}
	if lastBlockSize != 0 {
		res, err := encodeBlock(src[fullBlockCount*8:])
		if err != nil {
			return "", err
		}
		retVal += res
	}
	return retVal, nil
}

/*
DecodeFromString

This function takes a string and attempts to decode it using Monero Base58 formatting
*/
func DecodeFromString(src string) ([]byte, error) {
	lenData := len(src)
	if lenData == 0 {
		return nil, nil
	}
	fullBlockCount := lenData / encodedBlockSize
	lastBlockSize := lenData % encodedBlockSize
	var retVal []byte
	for i := 0; i < fullBlockCount; i++ {
		res, err := decodeBlock(src[i*encodedBlockSize : (i+1)*encodedBlockSize])
		if err != nil {
			return nil, err
		}
		retVal = append(retVal, res...)
	}
	if lastBlockSize != 0 {
		res, err := decodeBlock(src[fullBlockCount*11:])
		if err != nil {
			return nil, err
		}
		retVal = append(retVal, res...)
	}
	return retVal, nil
}
