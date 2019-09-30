package context

var bcContext BCContext

type BCContext interface {
	LoadContract(key []byte) ([]byte, error)
	Height() int64
}

func SetBCContext(context BCContext) {
	bcContext = context
}

func GetBCContext() BCContext {
	return bcContext
}

