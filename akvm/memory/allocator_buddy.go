package memory

import (
	"fmt"
	"math"

	ankrcmm "github.com/Ankr-network/ankr-chain/common"
	"github.com/go-interpreter/wagon/exec/common"
)

func leftChildIndex(index uint64) uint64 {
	return index*2 + 1
}

func rightChildIndex(index uint64) uint64 {
	return index*2 + 2
}

func parentIndex(index uint64) uint64 {
	return (index + 1)/2 -1
}

type allocatorBuddy struct {
	totalSize uint
	nodesSize []uint
}

func newAllocatorBuddy(totalSize uint) *allocatorBuddy {
	if totalSize < 1 || !common.IsPowOf2(totalSize) {
		return nil
	}

	allocBuddy := new(allocatorBuddy)
	allocBuddy.totalSize = totalSize
	allocBuddy.nodesSize = make([]uint, 2 * totalSize - 1)

	nodeSizeTemp := totalSize * 2

	for i := uint(0); i <  2 * totalSize - 1; i++ {
		if common.IsPowOf2(i+1) {
			nodeSizeTemp /= 2
		}

		allocBuddy.nodesSize[i] = nodeSizeTemp
	}

	return allocBuddy
}

func (allocator *allocatorBuddy) alloc(size uint) (uint64, error) {
    if size == 0 {
    	size = 1
	}

    if !common.IsPowOf2(size) {
    	size = common.FixSize(size)
	}

    if allocator.nodesSize[0] < size {
    	return 0, fmt.Errorf("the alloc size(%d) beyond total size(%d)", size, allocator.totalSize)
	}

    index        := uint64(0)
	nodeSizeTemp := uint(0)

    for nodeSizeTemp = allocator.totalSize; nodeSizeTemp != size; nodeSizeTemp /= 2 {
    	if allocator.nodesSize[leftChildIndex(index)] >= size{
    		index = leftChildIndex(index)
		}else {
			index = rightChildIndex(index)
		}
	}

    allocator.nodesSize[index] = 0
    offset := (index + 1) * uint64(nodeSizeTemp) - uint64(allocator.totalSize)

    for index > 0 {
    	index = parentIndex(index)
    	leftNodeSize  := allocator.nodesSize[leftChildIndex(index)]
    	rightNodeSize := allocator.nodesSize[rightChildIndex(index)]
    	allocator.nodesSize[index] = ankrcmm.MaxUint(leftNodeSize, rightNodeSize)
	}

    return offset, nil
}

func (allocator *allocatorBuddy) free(offset uint64) error {
	if offset < 0 || offset >= uint64(allocator.totalSize) {
		return fmt.Errorf("offset beyong memory limit: offset=%d, totalSize=%d", offset, allocator.totalSize)
	}

	nodeSizeTemp := uint(1)
	index        := offset + uint64(allocator.totalSize) -1

    for ; allocator.nodesSize[index] > 0; index = parentIndex(index) {
		nodeSizeTemp *= 2
		if index == 0 {
			return nil
		}
	}

    allocator.nodesSize[index] = nodeSizeTemp

    for index > 0 {
    	index = parentIndex(index)
    	nodeSizeTemp *= 2

		leftNodeSize  := allocator.nodesSize[leftChildIndex(index)]
		rightNodeSize := allocator.nodesSize[rightChildIndex(index)]
		if leftNodeSize + rightNodeSize == nodeSizeTemp {
			allocator.nodesSize[index] = nodeSizeTemp
		}else {
			allocator.nodesSize[index] = uint(math.Float64bits(math.Max(float64(leftNodeSize), float64(rightNodeSize))))
		}
	}

    return nil
}

func (allocator *allocatorBuddy) size(offset uint64) (uint, error) {
	if offset < 0 || offset >= uint64(allocator.totalSize) {
		return 0, fmt.Errorf("offset beyong memory limit: offset=%d, totalSize=%d", offset, allocator.totalSize)
	}

	nodeSizeTemp := uint(1)
	for index := offset  + uint64(allocator.totalSize-1); allocator.nodesSize[index] > 0; index = parentIndex(index) {
		nodeSizeTemp  *= 2
	}

	return nodeSizeTemp, nil
}

func (allocator *allocatorBuddy) growTotalSize(size uint) error {
	curTotalSize := allocator.totalSize

	groupBy := (size + curTotalSize)/curTotalSize
	if !common.IsPowOf2(groupBy) {
		return fmt.Errorf("invalid grow size: size=%d", size)
	}

	newTotalSize := size + curTotalSize
	newNodesSize := make([]uint, 2 * newTotalSize - 1)

	nodeSizeTemp := 2 * newTotalSize

	j    := uint(0)
	t    := uint(0)
	cntl := uint(1)
	for i := uint(0); i <  2 * newTotalSize - 1; i++ {
		if common.IsPowOf2(i+1) {
			nodeSizeTemp /= 2
			cntl *= 2
			t=0
		}

		if nodeSizeTemp <= curTotalSize && t < cntl/2 {
			newNodesSize[i] = allocator.nodesSize[j]
			j++
			t++
		} else {
			newNodesSize[i] = nodeSizeTemp
		}
	}

	allocator.totalSize = newTotalSize
	allocator.nodesSize = newNodesSize

	return nil
}


