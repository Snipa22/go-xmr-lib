package serialization

import "errors"

var ErrNonceTooLong = errors.New("nonce is too long to insert into the transaction data field")
