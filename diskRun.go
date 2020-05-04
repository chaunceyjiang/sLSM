package sLSM

import (
	"encoding/binary"
	"log"
	"os"
	"strconv"
	"syscall"
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
	fencePointersOffset []blockHandle
	data                []byte
	cmp                 Comparer
}

func NewDiskRun(capacity uint64, level int, runID int, bfFp float64, cmp Comparer) *DiskRun {
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
		cmp:           cmp,
	}
}

func (dr *DiskRun) writeData(pairs []KVPair) {
	var idxOffset uint64 = 0
	//dr.fd.Seek(offset, io.SeekStart)
	// 保存数据
	for i := 0; i < len(pairs); i++ {
		dr.fd.Write(append(pairs[i].Key, pairs[i].Value...))
		dr.fencePointers = append(dr.fencePointers, pairs[i].Key)
		dr.fencePointersOffset = append(dr.fencePointersOffset,
			blockHandle{offset: idxOffset, length: uint64(len(pairs[i].Key))})
		idxOffset += uint64(len(pairs[i].Key)) + uint64(len(pairs[i].Value))
	}
	// 保存索引
	for i := 0; i < len(dr.fencePointers); i++ {
		tmp := make([]byte, 2*binary.MaxVarintLen64)
		n := encodeBlockHandle(tmp, dr.fencePointersOffset[i])
		dr.fd.Write(tmp[:n])
	}
	// 保存索引的开始位置
	dr.fd.Write(uint642byte(idxOffset))
	// 刷盘
	if err := dr.fd.Sync(); err != nil {
		log.Println(err)
	}
	dr.closed = true
	if err := dr.fd.Close(); err != nil {
		log.Println(err)
	}
	dr.fencePointersOffset = dr.fencePointersOffset[:]
	dr.fencePointers = dr.fencePointers[:]
}

func (dr *DiskRun) Lookup(key K) V {
	var err error
	if dr.closed {
		if dr.fd, err = os.Open(dr.filename); err != nil {
			log.Fatalln(err)
		}
	}
	size, _ := dr.fd.Stat()
	dr.data, _ = syscall.Mmap(int(dr.fd.Fd()), 0, int(size.Size()), syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	idxOffset := byte2uint64(dr.data[len(dr.data)-8:])
	for int(idxOffset) < len(dr.data)-8 {
		bh, n := decodeBlockHandle(dr.data[idxOffset:])
		if dr.cmp.Eq(dr.data[bh.offset:bh.offset+bh.length], key) {
			return // TODO
		}
		//dr.fencePointers = append(dr.fencePointers, dr.data[bh.offset:bh.offset+bh.length])
		//dr.fencePointersOffset = append(dr.fencePointersOffset, bh)
		idxOffset += uint64(n)
	}
	return nil
}
