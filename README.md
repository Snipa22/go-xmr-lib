# go-xmr-lib
Golang library for supporting RPC communication with Monero based daemons and RPC wallets

# CNUtil
* address_decode - Complete, all address decoding steps are working properly
* convert_blob - Implemented as GetBlockHashingBlob - This function generates the actual hashable string from a block struct
* construct_block_blob - Implemented as ParseBlockFromTemplateBlob - This function generates the internal block struct from a block template
* get_block_id - TODO - Keccak256 over the block blob data.  Pack it all together then run the hash


The following are API calls, and likely need to be built as part of a larger RPC library, though this project is fine to contain them based on the face this is a general XMR lib, it may be worth extracting into it's own library in the future depending on changes to the RPC protocols.  
Should implement these as a common library going forwards with consideration, some are JSON-RPC, some are direct GET/POST calls, so that needs to be reviewed against the docs, linked in each subsection
For JsonRPC, use a two-stage decoder that can decode the json RPC error wrapper, then export the actual decoded object as a go struct.
#TODO - Daemon
* Docs - https://web.getmonero.org/resources/developer-guides/daemon-rpc.html
* getblockheaderbyhash - Primarily used in the payment loop, can wait for backend impl
* getblockheaderbyheight - Unused
* getlastblockheader - Used in blockmanager and Backend
* submitblock - BlockRepeater and main pool (submits blocks)
* getblocktemplate - Used in main pool loops to get the block `templates`

#TODO - Wallet
* Docs - https://www.getmonero.org/resources/developer-guides/wallet-rpc.html
* getbalance - Gets the balance of the wallet
* getheight - Gets the current height of the wallet, should be the same as the current highest BT
* store - Saves the wallet to disk in case of shutdown (Should be run on cron.)
* transfer - The actual transfer function that sends money around.