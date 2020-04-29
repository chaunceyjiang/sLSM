package sLSM

import (
	"fmt"
	"testing"
)

func TestMurMurHash128(t *testing.T) {
	hash128:=New128()
	hash128.Write([]byte("test"))
	fmt.Println(hash128.Sum128())
}