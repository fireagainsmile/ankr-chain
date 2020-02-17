package common

type ChainStateInfo interface {
	LatestHeight() int64
}
