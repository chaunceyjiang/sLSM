package sLSM

import "math"

// bloom表达式实现
// 公式
// 由 用户输入元素个数n + 误差率 p
// 计算可得 内存大小m bites ,k 个hash function
// m=-n*lnP/(ln2)^2
//
type BloomFilter struct {
	m    uint64
	k    uint64
	bits []bool
}

func NewBloomFilter(n uint64, p float64) *BloomFilter {
	m := -1 * (float64(n) * math.Log(p)) / (math.Ln2 * math.Ln2)
	k := math.Ceil((math.Ln2 * m) / float64(n))
	return &BloomFilter{
		m:    uint64(m),
		k:    uint64(k),
		bits: make([]bool, uint64(m)),
	}
}

func (bf *BloomFilter) hash(key K) (h1 uint64, h2 uint64) {
	return Sum128(key)
}

func (bf *BloomFilter) nthHash(n uint64, h1 uint64, h2 uint64) uint64 {
	return (h1 + n*h2) % bf.m
}

func (bf *BloomFilter) Add(key K) {
	h1, h2 := bf.hash(key)
	for i := uint64(0); i < bf.k; i++ {
		bf.bits[bf.nthHash(i, h1, h2)] = true
	}
}

func (bf *BloomFilter) MayContain(key K) bool {
	h1, h2 := bf.hash(key)
	for i := uint64(0); i < bf.k; i++ {
		// 只要有一个不在，就一定不在
		if !bf.bits[bf.nthHash(i, h1, h2)]{
			return false
		}
	}
	// false negative
	return true
}
