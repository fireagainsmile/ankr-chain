module github.com/Ankr-network/ankr-chain

go 1.12

replace github.com/tendermint/tendermint => github.com/Ankr-network/tendermint v0.31.5-0.20190915180111-98f9e8c1668c

replace github.com/go-interpreter/wagon => github.com/Ankr-network/wagon v0.9.0-0.20191015152132-a57bd86fecb0

require (
	github.com/Ankr-network/dccn-common v0.0.0-20191014090437-9fa44d3777fe
	github.com/VividCortex/gohistogram v1.0.0 // indirect
	github.com/btcsuite/btcd v0.20.0-beta // indirect
	github.com/fortytw2/leaktest v1.3.0 // indirect
	github.com/go-interpreter/wagon v0.6.0
	github.com/go-sql-driver/mysql v1.4.1
	github.com/rcrowley/go-metrics v0.0.0-20190826022208-cac0b30c2563 // indirect
	github.com/rs/cors v1.7.0 // indirect
	github.com/spf13/cobra v0.0.5
	github.com/spf13/viper v1.4.0
	github.com/stretchr/testify v1.4.0
	github.com/tendermint/go-amino v0.15.1
	github.com/tendermint/iavl v0.12.4
	github.com/tendermint/tendermint v0.32.6
	gopkg.in/mgo.v2 v2.0.0-20190816093944-a6b53ec6cb22
)
