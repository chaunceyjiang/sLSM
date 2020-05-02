package sLSM

import (
	"errors"
	"fmt"
	"math/rand"
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
	maxKey   K
	minKey   K
	curLevel int
	p        float64
	cmp      Comparer
	n        int64
}

func NewSkipList(cmp Comparer) *SkipList {
	return &SkipList{
		head:     newNode(nil, nil, MaxLevel),
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
	if sk.minKey == nil || sk.cmp.Lt(key, sk.minKey) {
		sk.minKey = key
	}
	if sk.maxKey == nil || sk.cmp.Gt(key, sk.maxKey) {
		sk.maxKey = key
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
			for i := sk.curLevel + 1; i <= lvl; i++ {
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

func (sk *SkipList) Search(key K) (V, bool) {
	if !checkKey(key) {
		return nil, false
	}
	cur := sk.head
	for i := sk.curLevel; i >= 0; i-- {
		for cur.forward[i] != nil && sk.cmp.Gt(key, cur.forward[i].key) {
			cur = cur.forward[i]
		}
	}
	cur = cur.forward[0]

	if cur == nil || sk.cmp.Neq(cur.key, key) {
		return nil, false
	}
	return cur.value, true
}

func (sk *SkipList) Delete(key K) {
	if !checkKey(key) {
		return
	}
	updated := make([]*node, MaxLevel+1)
	cur := sk.head
	for i := sk.curLevel; i >= 0; i-- {
		for cur.forward[i] != nil && sk.cmp.Gt(key, cur.forward[i].key) {
			cur = cur.forward[i]
		}
		updated[i] = cur
	}
	cur = cur.forward[0]

	if sk.cmp.Neq(cur.key, key) {
		return
	}
	// 修改每一层指针
	for i := 0; i < sk.curLevel; i++ {
		if updated[i].forward[i] != cur {
			break
		}
		// 修改删除节点的前一个指针
		updated[i].forward[i] = cur.forward[i]
	}
	sk.n--
	for sk.curLevel > 1 && sk.head.forward[sk.curLevel] == nil {
		// 头结点指向了一个空值，表示当前层已经没有节点了，层数减一
		sk.curLevel--
	}
}

// 产生一个随机层
func (sk *SkipList) randomLevel() int {
	lvl := 0
	for rand.Float64() < sk.p && lvl < MaxLevel {
		lvl++
	}
	return lvl
}

func (sk *SkipList) ShowList(f func(key K) string) {
	fmt.Println("-----skip list-----")
	for i := sk.curLevel; i >= 0; i-- {
		cur := sk.head.forward[i]
		for cur != nil {
			fmt.Printf(f(cur.key))
			if cur.forward[i] != nil {
				fmt.Printf("->")
			}
			cur = cur.forward[i]
		}
		fmt.Println()
	}
}

func (sk *SkipList) getAll() []KVPair {
	node := sk.head.forward[0]
	rest := make([]KVPair, 0)
	for node != nil {
		rest = append(rest, KVPair{
			Key:   node.key,
			Value: node.value,
		})
		node = node.forward[0]
	}
	return rest
}

func (sk *SkipList) getAllInRange(k1, k2 K) []KVPair {
	if sk.cmp.Gt(k1, sk.maxKey) || sk.cmp.Lt(k2, sk.minKey) {
		return nil
	}
	rest := make([]KVPair, 0)
	node := sk.head.forward[0]
	for sk.cmp.Lt(node.key, k1) {
		node = node.forward[0]
	}
	for sk.cmp.Lt(node.key, k2) {

		rest = append(rest, KVPair{
			Key:   node.key,
			Value: node.value,
		})
		node = node.forward[0]
	}

	return rest
}
