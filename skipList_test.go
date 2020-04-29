package sLSM

import (
	"encoding/binary"
	"testing"
)

func TestSkipList_Insert(t *testing.T) {
	sk := &SkipList{
		head:     newNode(nil, nil, MaxLevel),
		max:      MaxLevel,
		curLevel: 0,
		p:        0.5,
		cmp:      testCmp{},
		n:        0,
	}
	var count = 10
	var key ,value []byte
	for i := count; i < 4*count; i++ {
		key = make([]byte, 8)
		value = make([]byte, 8)
		binary.LittleEndian.PutUint64(key, uint64(i*2))
		binary.LittleEndian.PutUint64(value, uint64(i))
		_ = sk.Insert(key, value)
	}
	sk.showList()
	key = make([]byte, 8)
	value = make([]byte, 8)
	binary.LittleEndian.PutUint64(key, uint64(19))
	binary.LittleEndian.PutUint64(value, uint64(27))
	_ = sk.Insert(key, value)
	sk.showList()
	key1 := make([]byte, 8)
	value1 := make([]byte, 8)
	binary.LittleEndian.PutUint64(key1, uint64(27))
	binary.LittleEndian.PutUint64(value1, uint64(27))
	_ = sk.Insert(key1, value1)
	sk.showList()
}
