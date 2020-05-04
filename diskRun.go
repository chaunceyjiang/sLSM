package sLSM

import (
	"io"
	"log"
	"os"
	"strconv"
)

// 每个磁盘上的runs和内存中的runs类似，都是有max/min key以及Bloom filter进行过滤和fencePointer索引
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
	var f *os.File
	var err error
	if f, err = os.Create(filename); err != nil {
		log.Fatalln(err)
	}
	//if !fileExist(filename) {
	//	if f, err = os.Create(filename); err != nil {
	//		log.Fatalln(err)
	//	}
	//}
	//if f, err = os.OpenFile(filename, os.O_RDWR|os.O_APPEND|os.O_TRUNC, 0666); err != nil {
	//	log.Fatalln(err)
	//}

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
	idxOffset := 0
	//dr.fd.Seek(offset, io.SeekStart)
	// 保存数据
	for i := 0; i < len(pairs); i++ {
		dr.fd.Write(append(pairs[i].Key, pairs[i].Value...))
		dr.fencePointers = append(dr.fencePointers, pairs[i].Key)
		dr.fencePointersOffset = append(dr.fencePointersOffset, idxOffset)
		idxOffset += len(pairs[i].Key) + len(pairs[i].Value)
	}
	// 保存索引
	for i := 0; i < len(dr.fencePointers); i++ {
		dr.fd.Write(append(dr.fencePointers[i], uint642byte(uint64(dr.fencePointersOffset[i]))...))
	}
	// 保存索引的开始位置
	dr.fd.Write(uint642byte(uint64(idxOffset)))
	// 刷盘
	if err := dr.fd.Sync(); err != nil {
		log.Println(err)
	}
}

func (dr *DiskRun) Lookup(key K) V {
	idx, exist := dr.getIndex(key)
	if exist {
		return nil
	}
	return
}

func (dr *DiskRun) getIndex(key K) (int, bool) {
	return 0, true
}
