package memory

import (
	"errors"

	"github.com/Ankr-network/wagon/wasm"
)

type HeapMemoryImpl struct {
	allocator *allocatorBuddy
}

func NewHeapMemory( ) wasm.HeapMemory {
	return &HeapMemoryImpl{}
}

func (hm *HeapMemoryImpl) Init(totalSize uint) {
	hm.allocator = newAllocatorBuddy(totalSize)
}

func (hm *HeapMemoryImpl) Malloc(size uint) (uint64, error) {
	if hm.allocator == nil {
		return 0, errors.New("HeapMemory's allocator nil")
	}

	return hm.allocator.alloc(size)
}

func (hm *HeapMemoryImpl) Free(offset uint64) error {
	if hm.allocator == nil {
		return errors.New("HeapMemory's allocator nil")
	}

	return hm.allocator.free(offset)
}

func (hm *HeapMemoryImpl) GrowMemory(size uint) error {
	if hm.allocator == nil {
		return errors.New("HeapMemory's allocator nil")
	}

	return hm.allocator.growTotalSize(size)
}
