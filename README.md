# AnkrChain

AnkrChain currently supports two functions: metering and ankr coin transfer.
Metering used the key value store, and each data center should have an account.
Ankr coin transfer used address and balance for each account, and allows coin transfer.


## Private Key, Public Key and Account
In the tool and api_sample folder under abci/ankrchain, there is genkey example to generate private key, public key and address.


# APIs
We have APIs to visit the ankrchain.

GOLANG Version APIs in [ANKRCHAIN GOLANG APIS](https://github.com/Ankr-network/dccn-common/tree/develop/wallet)

GOLANG Version API samples in abci/ankrchain/api_sample

GOLANG APIs

- wallet.GenerateKeys() //address length is 46 bytes
- wallet.GetBalance()
- wallet.SendCoins()
- wallet.Sign()
- wallet.GetAddressByPublicKey()
- wallet.GetStake()
- wallet.GetHistorySend()
- wallet.GetHistoryReceive()
- wallet.SetMetering()

GOLANG Admin:

- wallet.SetStake()
- wallet.SetBalance()


JAVASCRIPT VERSION APIs in [ANKRCHAIN JAVASCRIPT APIS](https://github.com/Ankr-network/dccn-common/tree/develop/frontend/wallet-sdk)

JAVASCRIPS APIs

- gen_key  // address length is 46 bytes
- get_balance
- send_coin
- get_history


## CURL Examples

CURL examples are not encouraged to use, unless for testing purpose. Users are encouraged to use APIs and api samples.

Query Balance(it will show 'not exist' and the address is new):

```
curl  'localhost:26657/abci_query?data="bal:1111111111222222222233333333335555555555"'
```

Set Balance: (for test)

```
curl -s 'localhost:26657/broadcast_tx_commit?tx="set_bal=1111111111222222222233333333335555555555:150"'
```

Send Coins: (address1 send 600 coins to address2 as below)
Note: address1 needs has enough balance.

```
curl -s 'localhost:26657/broadcast_tx_commit?tx="trx_send=B508ED0D54597D516A680E7951F18CAD24C7EC9F:1234567890123456789012345678901234567890:600:(nonce):(sig):(pub_key)"'
```

Initial Accounts:

When creating blockchain, there are two ways to allocate initial coins in genesis

1) in genesis config file, Ankr blockchain will try to read from config file.
2) if not found in config file, we pre-allocate 100000 coins for address B508ED0D54597D516A680E7951F18CAD24C7EC9F for test purpose.

