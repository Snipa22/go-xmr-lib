# go-xmr-lib
Golang library for supporting RPC communication with Monero based daemons and RPC wallets

# TODO - CNUtil
* address_decode
  * Need to implement Base58 XMR decoding to get full address
  * Implement Keccak-1600 to perform address checksumming for validity
* address_decode_integrated
* convert_blob
* construct_block_blob
* get_block_id

#TODO - Daemon
* getblockheaderbyhash
* getblockheaderbyheight
* getlastblockheader
* submitblock
* getblocktemplate

#TODO - Wallet
* getbalance
* getheight
* store
* transfer