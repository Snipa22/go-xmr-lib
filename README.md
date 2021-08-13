# go-xmr-lib
Golang library for supporting RPC communication with Monero based daemons and RPC wallets

# CNUtil
* address_decode - Complete, all address decoding steps are working properly
  * Need to implement Base58 XMR decoding to get full address
  * Implement Keccak-1600 to perform address checksumming for validity
* address_decode_integrated - TODO
* convert_blob - TODO
* construct_block_blob - TODO
* get_block_id - TODO

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