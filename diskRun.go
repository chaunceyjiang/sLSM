package sLSM

import (
	"log"
	"os"
	"strconv"
)

// 每个磁盘上的runs和内存中的runs类似，都是有max/min key以及Bloom filter进行过滤和索引
type DiskRun struct {
	closed   bool
	filename string
	fd       *os.File
	capacity uint64

	minKey              K
	maxKey              K
	bf                  *BloomFilter
	fencePointers       []K
	fencePointersOffset []int
	data                []byte
}

func NewDiskRun(capacity uint64, level int, runID int, bfFp float64) *DiskRun {
	filename := "C_" + strconv.Itoa(level) + "_" + strconv.Itoa(runID) + ".txt"
	f, err := os.Create(filename)
	if err != nil {
		log.Fatalln(err)
	}
	return &DiskRun{
		closed:        false,
		filename:      filename,
		fd:            f,
		capacity:      capacity,
		minKey:        nil,
		maxKey:        nil,
		bf:            NewBloomFilter(capacity, bfFp),
		fencePointers: make([]K, capacity),
		data:          nil,
	}
}

func (dr *DiskRun) writeData(pairs []KVPair) {
	offset := 0
	for i := 0; i < len(pairs); i++ {
		dr.fd.Write(append(pairs[i].Key, pairs[i].Value...))
		dr.fencePointers = append(dr.fencePointers, pairs[i].Key)
		dr.fencePointersOffset = append(dr.fencePointersOffset, offset)
		offset += len(pairs[i].Key) + len(pairs[i].Value)
	}

}
