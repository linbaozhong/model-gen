package distance

import (
	"fmt"
	"testing"
)

func BenchmarkGoogle(b *testing.B) {
	for i := 0; i < b.N; i++ {
		google()
	}
}

func TestGoogle(t *testing.T) {
	fmt.Println("google: ", google())
}

func google() float64 {
	return Distance(116.368904, 39.923423, 116.387271, 39.922501)
}
