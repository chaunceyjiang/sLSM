package sLSM

import (
	"fmt"
	"testing"
)

func TestNewHashTable(t *testing.T) {
	ht := NewHashTable(5, testCmp{})
	for i := 0; i < 10; i++ {
		ht.Put(genTestKeyValue(i))
	}
	fmt.Println(ht.Get(genTestKeyValue(1)))
	fmt.Println(ht.Get(genTestKeyValue(5)))
	fmt.Println(ht.Get(genTestKeyValue(9)))
}
