module github.com/Ankr-network/ankr-chain

go 1.12

replace github.com/tendermint/tendermint => github.com/Ankr-network/tendermint v0.31.5-0.20191016091852-c60735e225bb

replace github.com/go-interpreter/wagon => github.com/Ankr-network/wagon v0.9.0-0.20191015152132-a57bd86fecb0

require (
	github.com/Ankr-network/dccn-common v0.0.0-20191014090437-9fa44d3777fe
	github.com/VividCortex/gohistogram v1.0.0 // indirect
	github.com/btcsuite/btcd v0.20.0-beta // indirect
	github.com/fortytw2/leaktest v1.3.0 // indirect
	github.com/go-interpreter/wagon v0.6.0
	github.com/go-sql-driver/mysql v1.4.1
	github.com/jmhodges/levigo v1.0.0 // indirect
	github.com/rcrowley/go-metrics v0.0.0-20190826022208-cac0b30c2563 // indirect
	github.com/rs/cors v1.7.0 // indirect
	github.com/spf13/cobra v0.0.5
	github.com/spf13/viper v1.4.0
	github.com/stretchr/testify v1.4.0
	github.com/syndtr/goleveldb v1.0.0 // indirect
	github.com/tendermint/go-amino v0.15.0
	github.com/tendermint/iavl v0.12.0
	github.com/tendermint/tendermint v0.32.6
	google.golang.org/genproto v0.0.0-20180831171423-11092d34479b // indirect
	gopkg.in/mgo.v2 v2.0.0-20190816093944-a6b53ec6cb22
)
