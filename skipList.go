package sLSM

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
)

const MaxLevel = 10

// 跳跃表实现
type node struct {
	key     K
	value   V
	forward []*node
}

func newNode(key K, value V, level int) *node {
	return &node{
		key:     key,
		value:   value,
		forward: make([]*node, level+1),
	}
}

type SkipList struct {
	head     *node
	max      int
	curLevel int
	p        float64
	cmp      Comparer
	n        int64
}

func NewSkipList(cmp Comparer) *SkipList {
	return &SkipList{
		head:     newNode(nil, nil, MaxLevel),
		max:      MaxLevel,
		curLevel: 0,
		p:        0.5,
		cmp:      cmp,
		n:        0,
	}
}

func (sk *SkipList) Insert(key K, value V) error {
	if !checkKV(key, value) {
		return errors.New("key or value error")
	}
	updated := make([]*node, MaxLevel+1)
	cur := sk.head
	for i := sk.curLevel; i >= 0; i-- {
		for cur.forward[i] != nil && sk.cmp.Gt(key, cur.forward[i].key) {
			cur = cur.forward[i]
		}
		updated[i] = cur // 记录路径
	}

	cur = cur.forward[0]

	if cur == nil || sk.cmp.Neq(cur.key, key) {
		lvl := sk.randomLevel()
		if lvl > sk.curLevel {
			// 这里 i = curLevel +  1
			for i := sk.curLevel+1; i <= lvl; i++ {
				updated[i] = sk.head
			}
			sk.curLevel = lvl
		}

		newNode := newNode(key, value, lvl)

		// 这里有 i等于lvl
		for i := 0; i <= lvl; i++ {
			// 插入
			newNode.forward[i] = updated[i].forward[i]
			updated[i].forward[i] = newNode
		}
		sk.n++
	} else {
		// 更新
		cur.value = value
	}

	return nil
}

// 产生一个随机层
func (sk SkipList) randomLevel() int {
	lvl := 0
	for rand.Float64() < sk.p && lvl < MaxLevel {
		lvl++
	}
	return lvl
}

func (sk SkipList) showList() {
	fmt.Println("-----skip list-----")
	for i := sk.curLevel; i >= 0; i-- {
		cur := sk.head.forward[i]
		for cur != nil {
			fmt.Printf(strconv.Itoa(int(binary.LittleEndian.Uint64(cur.key))))
			if cur.forward[i] != nil {
				fmt.Printf("->")
			}
			cur = cur.forward[i]
		}
		fmt.Println()
	}
}
