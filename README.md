# go-xmr-lib
Golang library for supporting RPC communication with Monero based daemons and RPC wallets

# CNUtil
* address_decode - Complete, all address decoding steps are working properly
  * Need to implement Base58 XMR decoding to get full address
  * Implement Keccak-1600 to perform address checksumming for validity
* convert_blob - Implemented as GetBlockHashingBlob - This function generates the actual hashable string from a block struct
* construct_block_blob - Implemented as ParseBlockFromTemplateBlob - This function generates the internal block struct from a block template
* get_block_id - TODO - Keccak256 over the block blob data.  Pack it all together then run the hash

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