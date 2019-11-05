package memory

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMultiMalloc(t *testing.T) {
	heapM :=  NewHeapMemory( )
	assert.NotEqual(t, heapM, nil)
	heapM.Init(256)

	index1, _ := heapM.Malloc(10)
	index2, _ := heapM.Malloc(20)

	assert.NotEqual(t, index1, index2)
}