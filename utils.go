package sLSM

import "encoding/binary"

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
	return binary.LittleEndian.Uint64(k1) <binary.LittleEndian.Uint64(k2)
}

func (t testCmp) Eq(k1 K, k2 K) bool {
	return binary.LittleEndian.Uint64(k1) == binary.LittleEndian.Uint64(k2)
}

func (t testCmp) Neq(k1 K, k2 K) bool {
	return binary.LittleEndian.Uint64(k1) != binary.LittleEndian.Uint64(k2)
}