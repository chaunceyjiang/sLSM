package sLSM

// 开放地址 一次探测法
type HashTable struct {
	cmp     Comparer
	size    uint64
	curSize uint64
	table   []*KVPair
}

func NewHashTable(size uint64, cmp Comparer) *HashTable {
	table := make([]*KVPair, size*2)
	var i uint64
	for i = 0; i < size*2; i++ {
		table[i] = EMPTY
	}
	return &HashTable{
		cmp:     cmp,
		size:    size * 2,
		curSize: 0,
		table:   table,
	}
}
func (ht *HashTable) hashFunc(key K) uint64 {
	h1, _ := Sum128(key)
	return h1 % ht.size
}

func (ht *HashTable) resize() {
	ht.size *= 2
	newTable := make([]*KVPair, ht.size)
	var i, j uint64
	// 重新计算位置
	for i = 0; i < ht.size; i++ {
		newTable[i] = EMPTY
	}

	for i = 0; i < ht.size/2; i++ {
		if ht.table[i] != EMPTY {
			newHash := ht.hashFunc(ht.table[i].Key)
			for j = 0; ; j++ {
				if newTable[(newHash+j)%ht.size] == EMPTY {
					newTable[(newHash+j)%ht.size] = ht.table[i]
					break
				}
			}
		}
	}
	ht.table = newTable
}

func (ht *HashTable) Put(key K, value V) {
	if ht.curSize*2 > ht.size {
		ht.resize()
	}
	hashKey := ht.hashFunc(key)

	var i uint64
	for i = 0; ; i++ {
		if ht.table[(hashKey+i)%ht.size] == EMPTY {
			ht.table[(hashKey+i)%ht.size] = &KVPair{
				Key:   key,
				Value: value,
			}
			ht.curSize++
			return
		} else if ht.cmp.Eq(ht.table[(hashKey+i)%ht.size].Key, key) {
			ht.table[(hashKey+i)%ht.size].Value = value
			return
		}
	}
}
func (ht *HashTable) Get(key K, value V) (V, bool) {
	hashKey := ht.hashFunc(key)
	var i uint64
	for i = 0; ; i++ {
		if ht.table[(hashKey+i)%ht.size] == EMPTY {
			return value, false
		} else if ht.cmp.Eq(ht.table[(hashKey+i)%ht.size].Key, key) {
			return ht.table[(hashKey+i)%ht.size].Value, true
		}
	}
}
