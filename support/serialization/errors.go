package serialization

import "errors"

var ErrNonceTooLong = errors.New("nonce is too long to insert into the transaction data field")
var ErrUnhandledTxExtraTag = errors.New("tx_extra contains an unrecognized tag byte")
var ErrTxExtraTruncated = errors.New("tx_extra data is truncated and shorter than the tag requires")
