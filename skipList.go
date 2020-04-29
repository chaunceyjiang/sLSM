package sLSM

const MAXLEVEL = 10

// 跳跃表实现
type node struct {
	key     []byte
	value   []byte
	forward []*node
}

func newNode(key []byte, value []byte, level int) *node {
	return &node{
		key:     key,
		value:   value,
		forward: make([]*node, level+1),
	}
}

type SkipList struct {
	head *node
	tail *node
	max int
	curLevel int
	p float64
}

func NewSkipList() *SkipList {
	return &SkipList{
		head:     nil,
		tail:     nil,
		max:      MAXLEVEL,
		curLevel: 0,
		p:        0.5,
	}
}

func (sk *SkipList) Insert(key []byte, value []byte) {
	
}

