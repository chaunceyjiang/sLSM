package sLSM

import (
	"testing"
)

func TestBloomFilter(t *testing.T) {
	bf := NewBloomFilter(100, 0.001)
	for i := 1; i <= 10; i++ {
		key1, _ := genTestKeyValue(i)
		bf.Add(key1)
	}
	key1, _ := genTestKeyValue(1)
	if !bf.MayContain(key1) {
		t.FailNow()
	}
	key1, _ = genTestKeyValue(11)
	if bf.MayContain(key1) {
		t.FailNow()
	}
}
