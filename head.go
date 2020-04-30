package sLSM

type KVIntPair struct {
	kvp KVPair
	i   int
}

func NewKVIntPair(pair KVPair, i int) *KVIntPair {
	return &KVIntPair{
		kvp: pair,
		i:   i,
	}
}

type StaticHeap struct {
	h   []KVIntPair
	cmp Comparer
}

func NewStaticHeap(size int, cmp Comparer) *StaticHeap {
	return &StaticHeap{
		h:   make([]KVIntPair, size),
		cmp: cmp,
	}
}
func (s *StaticHeap) Len() int {
	return len(s.h)
}

func (s *StaticHeap) Less(i, j int) bool {
	if s.cmp.Lt(s.h[i].kvp.Key, s.h[j].kvp.Key) {
		return true
	}
	if s.cmp.Eq(s.h[i].kvp.Key, s.h[j].kvp.Key) && s.h[i].i < s.h[j].i {
		return true
	}
	return false
}

func (s *StaticHeap) Swap(i, j int) {
	s.h[i], s.h[j] = s.h[j], s.h[i]
	//(*s).h[i], (*s).h[j] = (*s).h[j], (*s).h[i]
}

func (s *StaticHeap) Push(x interface{}) {
	s.h = append(s.h, x.(KVIntPair))
}

func (s *StaticHeap) Pop() (v interface{}) {
	s.h, v = s.h[:s.Len()-1], s.h[s.Len()-1]
	return
}
