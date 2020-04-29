package sLSM

import (
	"testing"
)

func testInit() *SkipList {
	sk := &SkipList{
		head:     newNode(nil, nil, MaxLevel),
		max:      MaxLevel,
		curLevel: 0,
		p:        0.5,
		cmp:      testCmp{},
		n:        0,
	}
	var count = 10

	for i := count; i < 4*count; i++ {

		_ = sk.Insert(genTestKeyValue(i))
	}
	sk.ShowList(showLevelFunc())
	return sk
}

func TestSkipList_Insert(t *testing.T) {
	sk := testInit()

	_ = sk.Insert(genTestKeyValue(19))
	sk.ShowList(showLevelFunc())

	_ = sk.Insert(genTestKeyValue(27))
	sk.ShowList(showLevelFunc())
}
