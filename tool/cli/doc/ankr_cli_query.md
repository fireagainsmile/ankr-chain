## ankr_cli query

query information from ankr chain
### Sub commands

```
  PS D:\> ankr-cli query
  query information from ankr chain
  
  Usage:
    ankr_cli query [command]
  
  Available Commands:
    account            query account info
    balance            get the balance of an address.
    block              Get block at a given height. If no height is provided, it will fetch the latest block.
    consensusstate     ConsensusState returns a concise summary of the consensus state
    dumpconsensusstate dumps consensus state
    genesis            Get genesis file.
    numunconfirmedtxs  Get number of unconfirmed transactions.
    status             Get Ankr status including node info, pubkey, latest block hash, app hash, block height and time.
    transaction        transaction allows you to query the transaction results.
    unconfirmedtxs     Get unconfirmed transactions (maximum ?limit entries) including their number
    validators         Get the validator set at the given block height. If no height is provided, it will fetch the current validator set.
  
  Flags:
    -h, --help         help for query
        --nodeurl string   validator url
  
  Use "ankr_cli query [command] --help" for more information about a command.
```

### usage  
    global options 
        --nodeurl string       url of a validator 
    * balance query target account balance    
        options:
            -a, --address string <requried> the address of the target account
            --symbol string    token symbol (default "ANKR")
            
    * account query target account info    
        options:
            --address string <requried> the address of the target account
    * block,  Get block at a given height. If no height is provided, it will fetch the latest block. And you can use "detail" to show more information about transactions contained in block.  
        options: 
            --height int   height interval of the blocks to query. integer or block interval formatted as [from:to] are accepted 
    * consensusstate, get the summary of the consensus state.   
        options: 
            NULL
    * dumpconsensusstate,  dumps consensus state. 
        options: 
            NULL
    * genesis,  Get genesis file. 
        options: 
            NULL
    * numunconfirmedtxs,  Get number of unconfirmed transactions. 
        options: 
            NULL
    * status,  Get Ankr status including node info, pubkey, latest block hash, app hash, block height and time. 
        options: 
            NULL
    * transaction,  transaction allows you to query the transaction results with multiple conditions. 
        options:
            --approve    bool      Include a proof of the transaction inclusion in the block
            --txid       string    The transaction hash
            --creator    string    app creator
            --from       string    the from address contained in a transaction
            --height     string    block height. Input can be an exactly block height  or a height interval separate by ":", and height interval should be enclosed with "[]" or "()" which is mathematically open interval and close interval.
            --metering   string    query metering transaction, both datacenter name and namespace should be  provided and separated  by ":"
            --page       int       Page number (1 based) (default 1)
            --perpage    int       Number of entries per page(max: 100) (default 30)
            --timestamp  string    transaction executed timestamp. Input can be an exactly unix timestamp  or a time interval separate by ":", and time interval should be enclosed with "[]" or "()" which is mathematically open interval and close interval.
            --to         string    the to address contained in a transaction
            --txid       string    The transaction hash
            --type       string    Ankr chain predefined types, SetMetering, SetBalance, UpdatValidator, SetStake, Send
    * unconfirmedtxs,  unconfirmed transactions including their number.
        options: 
            --limit int   number of entries (default 30)
    * validators,  Get the validator set at the given block heigh. 
        options:
            --height int   block height 
### example  
+ balance    
    ```
     PS D:\> ankr-cli query balance --nodeurl http://localhost:26657 --address E1403CA0DC201F377E820CFA62117A48D4D612400C20D3
     {
             "amount": "10000000000000000000000000000"
     }
    ```    
+ account   
    ```
    PS D:\> ankr-cli query balance --nodeurl http://localhost:26657 --address E1403CA0DC201F377E820CFA62117A48D4D612400C20D3
    {
               "accounttype": "7",
               "nonce": "16",
               "address": "E1403CA0DC201F377E820CFA62117A48D4D6124",
               "pubkey": "",
               "asserts": [
                       {
                               "currency": {
                                       "symbol": "ANKR",
                                       "decimal": "18"
                               },
                               "value": "IE/OXj4lAmEQAAAA"
                       }
               ]
       }
    ```
+ block     
    ``` 
    PS D:\> ankr-cli query block --nodeurl http://localhost:26657 --height 631 
    Total count: 1
    
    Block info:
    Version:  {10 1}
    Chain-Id: test-chain-4J6Kvw
    Height:  48
    Time: 2019-11-04 10:21:05.3357335 +0000 UTC
    Number-Txs:  1
    Total-Txs: 1
    Last-block-id:  B6195927CDB7BDA8530BBB2FFB8B4C1B49794B33EB54A776B03E78FE7335338C:1:60B7C6EA3085
    Last-commit-hash: 846C4ED895D0F18C978BB9CEB18216EE27839C2CDE5C94D6C1FE60614A564D4B
    Data-hash:  68BE7B0F901EA0158DAA11ED31D28723867E6CC4A19A64FF424FD5CF03AA523D
    Validator: 1805057E3381391B344B55017820DE2FC49E42388B2FF7213AC59885477885B0
    Consensus:  048091BC7DDC283F77BFBF91D73C44DA58C3DF8A9CBC867405D8B7F3DAADA22F
    Version:  {10 1}
    App-hash: 0000000000000000
    Proposer-Address: F840E36AD0CF3564DD9EE01E2E5948A6F42ABE80B6D16B
    
    Transactions contained in block:
    {
            "type": "ankr-chain/tx/txMsg",
            "value": {
                    "chainid": "test-chain-4J6Kvw",
                    "nonce": "0",
                    "gaslimit": "A+g=",
                    "gasprice": {
                            "currency": {
                                    "symbol": "ANKR",
                                    "decimal": "18"
                            },
                            "value": "AWNFeF2KAAA="
                    },
                    "gasused": null,
                    "signs": [
                            {
                                    "pubkey": {
                                            "type": "tendermint/PubKeyEd25519",
                                            "value": "wvHG3EddBbXQHcyJal0CS/YQcNYtEbFYxejnqf9OhM4="
                                    },
                                    "signed": "/IAa8vAGF2sFcbq6ZJ6xY3u5HYZyLhiqFyFRcXzXOt3qKJTZwtMyQnjw1Q7yFfqpxmcuWPQZvddT+UTBgpyTAw==",
                                    "R": "",
                                    "S": "",
                                    "PubPEM": ""
                            }
                    ],
                    "memo": "test transfer",
                    "version": "1.0",
                    "data": {
                            "type": "ankr-chain/tx/token/tranferTxMsg",
                            "value": {
                                    "fromaddr": "B508ED0D54597D516A680E7951F18CAD24C7EC9FCFCD67",
                                    "toaddr": "065E37B3FC243B9FABB1519AB876E7632C510DC9324031",
                                    "amounts": [
                                            {
                                                    "currency": {
                                                            "symbol": "ANKR",
                                                            "decimal": "18"
                                                    },
                                                    "value": "U0RINexYAAA="
                                            }
                                    ]
                            }
                    }
            }
    }

    ```
+ consensusstate     
    ``` 
    PS D:\> ankr-cli query consensusstate --nodeurl http://localhost:26657
    {
        "round_state": {
            "height/round/step": "117983/0/1",
            "start_time": "2019-08-01T05:35:21.2209774Z",
            "proposal_block_hash": "",
            "locked_block_hash": "",
            "valid_block_hash": "",
            "height_vote_set": [
                {
                    "round": "0",
                    "prevotes": [
                        "nil-Vote"
                    ],
                    "prevotes_bit_array": "BA{1:_} 0/10 = 0.00",
                    "precommits": [
                        "nil-Vote"
                    ],
                    "precommits_bit_array": "BA{1:_} 0/10 = 0.00"
                }
            ]
        }
    }
    ```
+ dumpconsensusstate     
    ``` 
    PS D:\> ankr-cli query consensusstate --nodeurl http://localhost:26657
    {
        "round_state": {
            "height/round/step": "118015/0/1",
            "start_time": "2019-08-01T05:35:54.1263101Z",
            "proposal_block_hash": "",
            "locked_block_hash": "",
            "valid_block_hash": "",
            "height_vote_set": [
                {
                    "round": "0",
                    "prevotes": [
                        "nil-Vote"
                    ],
                    "prevotes_bit_array": "BA{1:_} 0/10 = 0.00",
                    "precommits": [
                        "nil-Vote"
                    ],
                    "precommits_bit_array": "BA{1:_} 0/10 = 0.00"
                }
            ]
        }
    }
    ```
+ genesis     
    ``` 
    PS D:\> ankr-cli query genesis --nodeurl http://localhost:26657
    {
        Genesis:{
            "genesis_time": 2019-07-24 10:44:03.9174995 +0000 UTC
            "chain_id": test-chain-0bOrck
            "consensus_params":{
            "block": {
                "max_bytes": 22020096,
                "max_gas": -1,
                "time_iota_ms": 1000
            },
            "evidence": {
                "max_age": 100000
            },
            "validator": {
                "pub_key_types": [
                    "ed25519"
                ]
            }
        }
        "validators": [
            {
                 "address": "6BAE4972E4DEA417461BFC8C41198E9F9DF5D8F13CE359",
                 "pub_key": {
                     "type": "tendermint/PubKeyEd25519",
                     "value": "voN51o4hX95tLWfF0FYg+JPEjzc71G9OCrje3k+IZ/c="
                  },
                 "power": "10",
                 "name": ""
            }
        ],
        "app_hash": ""
        }
     }
    ```
+ numunconfirmedtxs     
    ``` 
    PS D:\> ankr-cli query numunconfirmedtxs --nodeurl http://localhost:26657
    {
            "n_txs": "0",
            "total": "0",
            "total_bytes": "0",
            "txs": null
    }
    ```
+ status     
    ``` 
    PS D:\> ankr-cli query status --nodeurl http://localhost:26657
    {
            "node_info": {
                    "protocol_version": {
                            "p2p": "7",
                            "block": "10",
                            "app": "1"
                    },
                    "id": "2df293afa2a0ef722f9c73ea05bbdea28b1c8391503996",
                    "listen_addr": "tcp://0.0.0.0:26656",
                    "network": "test-chain-4J6Kvw",
                    "version": "0.31.5",
                    "channels": "4020212223303800",
                    "moniker": "DESKTOP-NF0AS58",
                    "other": {
                            "tx_index": "on",
                            "rpc_address": "tcp://0.0.0.0:26657"
                    }
            },
            "sync_info": {
                    "latest_block_hash": "250485A6C75A6AFC1923AF4DA9E8749F63DCEF5E8427BB3B996278C04E10C424",
                    "latest_app_hash": "0200000000000000",
                    "latest_block_height": "802",
                    "latest_block_time": "2019-11-04T10:33:58.9193676Z",
                    "catching_up": false
            },
            "validator_info": {
                    "address": "F840E36AD0CF3564DD9EE01E2E5948A6F42ABE80B6D16B",
                    "pub_key": {
                            "type": "tendermint/PubKeyEd25519",
                            "value": "ojlpHiDda9iGTKUDalu5dqtIvdWitBMRmDWq3J76Hxs="
                    },
                    "voting_power": "10"
            }
    }
    ```
+ transaction     
    ``` 
    PS D:\> ankr-cli query transaction --nodeurl http://localhost:26657 --txid 0x72fb3fa4735e2de3e56ab50a5d2ddcdbd019012b34a226dce0b7a3d2e13bddeb
    {
            "hash": "A6C50753C6CB41A0495C9BDB7E7635368F4C675B96A2F1AA353704CDC05C9B00",
            "height": 48,
            "index": 0,
            "tx_result": {
                    "gasWanted": "1000",
                    "gasUsed": "20",
                    "tags": [
                            {
                                    "key": "YXBwLmZyb21hZGRyZXNz",
                                    "value": "QjUwOEVEMEQ1NDU5N0Q1MTZBNjgwRTc5NTFGMThDQUQyNEM3RUM5RkNGQ0Q2Nw=="
                            },
                            {
                                    "key": "YXBwLnRvYWRkcmVzcw==",
                                    "value": "MDY1RTM3QjNGQzI0M0I5RkFCQjE1MTlBQjg3NkU3NjMyQzUxMERDOTMyNDAzMQ=="
                            },
                            {
                                    "key": "YXBwLnRpbWVzdGFtcA==",
                                    "value": "MTU3Mjg2Mjg2NjM3Mzg5NzgwMA=="
                            },
                            {
                                    "key": "YXBwLnR5cGU=",
                                    "value": "VHJhbnNmZXI="
                            }
                    ]
            },
            "tx": "rAK/ms8MChF0ZXN0LWNoYWluLTRKNkt2dxoCA+giFAoICgRBTktSEBISCAFjRXhdigAAMmkKJRYk3mQgwvHG3EddBbXQHcyJal0CS/YQcNYtEbFYxejnqf9OhM4SQPyAGvLwBhdrBXG6umSesWN7uR2Gci4YqhchUXF81zrd6iiU2cLTMkJ48NUO8hX6qcZnLlj0Gb3XU/lEwYKckwM6DXRlc3QgdHJhbnNmZXJCAzEuMEp6Seac0wouQjUwOEVEMEQ1NDU5N0Q1MTZBNjgwRTc5NTFGMThDQUQyNEM3RUM5RkNGQ0Q2NxIuMDY1RTM3QjNGQzI0M0I5RkFCQjE1MTlBQjg3NkU3NjMyQzUxMERDOTMyNDAzMRoUCggKBEFOS1IQEhIIU0RINexYAAA=",
            "proof": {
                    "RootHash": "",
                    "Data": null,
                    "Proof": {
                            "total": 0,
                            "index": 0,
                            "leaf_hash": null,
                            "aunts": null
                    }
            }
    }

    ```
    ```
    PS D:\> ankr-cli query transaction --to 065E37B3FC243B9FABB1519AB876E7632C510DC9324031 --nodeurl http://localhost:26657
    Total Tx Count: 1
    Transactions search result:
    {
            "hash": "A6C50753C6CB41A0495C9BDB7E7635368F4C675B96A2F1AA353704CDC05C9B00",
            "height": 48,
            "index": 0,
            "tx_result": {
                    "gasWanted": "1000",
                    "gasUsed": "20",
                    "tags": [
                            {
                                    "key": "YXBwLmZyb21hZGRyZXNz",
                                    "value": "QjUwOEVEMEQ1NDU5N0Q1MTZBNjgwRTc5NTFGMThDQUQyNEM3RUM5RkNGQ0Q2Nw=="
                            },
                            {
                                    "key": "YXBwLnRvYWRkcmVzcw==",
                                    "value": "MDY1RTM3QjNGQzI0M0I5RkFCQjE1MTlBQjg3NkU3NjMyQzUxMERDOTMyNDAzMQ=="
                            },
                            {
                                    "key": "YXBwLnRpbWVzdGFtcA==",
                                    "value": "MTU3Mjg2Mjg2NjM3Mzg5NzgwMA=="
                            },
                            {
                                    "key": "YXBwLnR5cGU=",
                                    "value": "VHJhbnNmZXI="
                            }
                    ]
            },
            "tx": "rAK/ms8MChF0ZXN0LWNoYWluLTRKNkt2dxoCA+giFAoICgRBTktSEBISCAFjRXhdigAAMmkKJRYk3mQgwvHG3EddBbXQHcyJal0CS/YQcNYtEbFYxejnqf9OhM4SQPyAGvLwBhdrBXG6umSesWN7uR2Gci4YqhchUXF81zrd6iiU2cLTMkJ48NUO8hX6qcZnLlj0Gb3XU/lEwYKckwM6DXRlc3QgdHJhbnNmZXJCAzEuMEp6Seac0wouQjUwOEVEMEQ1NDU5N0Q1MTZBNjgwRTc5NTFGMThDQUQyNEM3RUM5RkNGQ0Q2NxIuMDY1RTM3QjNGQzI0M0I5RkFCQjE1MTlBQjg3NkU3NjMyQzUxMERDOTMyNDAzMRoUCggKBEFOS1IQEhIIU0RINexYAAA=",
            "proof": {
                    "RootHash": "",
                    "Data": null,
                    "Proof": {
                            "total": 0,
                            "index": 0,
                            "leaf_hash": null,
                            "aunts": null
                    }
            }
    }
    ```
+ unconfirmedtxs     
    ``` 
   PS D:\> ankr-cli query unconfirmedtxs --nodeurl http://localhost:26657
   n_tx: 0
   total: 0
   total_bytes: 0
   transactions:
   []
    ```
+ validators
    ```
     PS D:\> ankr-cli query validators --nodeurl http://localhost:26657
     {
             "block_height": "1022",
             "validators": [
                     {
                             "address": "F840E36AD0CF3564DD9EE01E2E5948A6F42ABE80B6D16B",
                             "pub_key": {
                                     "type": "tendermint/PubKeyEd25519",
                                     "value": "ojlpHiDda9iGTKUDalu5dqtIvdWitBMRmDWq3J76Hxs="
                             },
                             "voting_power": "10",
                             "proposer_priority": "0"
                     }
             ]
     }
    ```
