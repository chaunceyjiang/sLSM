package sLSM

// 比较器
type Comparer interface {
	Gt(key []byte, value []byte) bool  // 大于
	Lt(key []byte, value []byte) bool  // 小于
	Eq(key []byte, value []byte) bool  // 等于
	Neq(key []byte, value []byte) bool //不等于
}

func checkKey(key []byte) bool {
	if len(key) == 0 {
		return false
	}
	return true
}

func checkVale(value []byte) bool {
	if len(value) == 0 {
		return false
	}
	return true
}
func checkKV(key, value []byte) bool {
	if checkKey(key) && checkVale(value) {
		return true
	}
	return false
}
