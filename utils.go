package sLSM

import (
	"encoding/binary"
	"os"
	"strconv"
)

// 比较器
type Comparer interface {
	Gt(k1 K, k2 K) bool  // 大于
	Lt(k1 K, k2 K) bool  // 小于
	Eq(k1 K, k2 K) bool  // 等于
	Neq(k1 K, k2 K) bool //不等于
}

type K []byte
type V []byte

type KVPair struct {
	Key   K
	Value V
}

var EMPTY = &KVPair{Key: nil, Value: nil}

func checkKey(key K) bool {
	if len(key) == 0 {
		return false
	}
	return true
}

func checkVale(value V) bool {
	if len(value) == 0 {
		return false
	}
	return true
}
func checkKV(key K, value V) bool {
	if checkKey(key) && checkVale(value) {
		return true
	}
	return false
}

type testCmp struct {
}

func (t testCmp) Gt(k1 K, k2 K) bool {
	return binary.LittleEndian.Uint64(k1) > binary.LittleEndian.Uint64(k2)
}

func (t testCmp) Lt(k1 K, k2 K) bool {
	return binary.LittleEndian.Uint64(k1) < binary.LittleEndian.Uint64(k2)
}

func (t testCmp) Eq(k1 K, k2 K) bool {
	return binary.LittleEndian.Uint64(k1) == binary.LittleEndian.Uint64(k2)
}

func (t testCmp) Neq(k1 K, k2 K) bool {
	return binary.LittleEndian.Uint64(k1) != binary.LittleEndian.Uint64(k2)
}

func showLevelFunc() func(key K) string {
	return func(key K) string {
		return strconv.Itoa(int(binary.LittleEndian.Uint64(key)))
	}
}

func genTestKeyValue(i int) (key K, value V) {
	key = make(K, 8)
	value = make(V, 8)
	binary.LittleEndian.PutUint64(key, uint64(i))
	binary.LittleEndian.PutUint64(value, uint64(i*2))
	return
}

func fileExist(filename string) bool {
	_, err := os.Stat(filename)
	if err != nil && os.IsNotExist(err){
		return false
	}
	return true
}
