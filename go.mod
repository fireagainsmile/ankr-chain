module github.com/Ankr-network/ankr-chain

go 1.12

replace github.com/tendermint/tendermint => github.com/Ankr-network/tendermint v0.31.5-0.20190719093344-1f8077fcd482

require (
	github.com/btcsuite/btcd v0.0.0-20190629003639-c26ffa870fd8 // indirect
	github.com/go-sql-driver/mysql v1.4.1
	github.com/rcrowley/go-metrics v0.0.0-20190706150252-9beb055b7962 // indirect
	github.com/rs/cors v1.6.0 // indirect
	github.com/spf13/cobra v0.0.5
	github.com/spf13/viper v1.4.0
	github.com/syndtr/goleveldb v1.0.0 // indirect
	github.com/tendermint/go-amino v0.15.0 // indirect
	github.com/tendermint/tendermint v0.0.0-00010101000000-000000000000
	google.golang.org/genproto v0.0.0-20180831171423-11092d34479b // indirect
	gopkg.in/mgo.v2 v2.0.0-20180705113604-9856a29383ce
)
