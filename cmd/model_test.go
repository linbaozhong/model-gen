package cmd

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

type Build struct {
	ID      int
	Name    string
	Version string
}

var buildPool = sync.Pool{
	New: func() interface{} {
		return new(Build)
	},
}

func NewBuild() *Build {
	return buildPool.Get().(*Build)
}

func (b *Build) Free() {
	b.ID = 0
	b.Name = ""
	b.Version = ""
	buildPool.Put(b)
}

func TestName(t *testing.T) {
	for i := 0; i < 10; i++ {
		go func() {
			obj := NewBuild()
			fmt.Printf("%v", obj)
			if obj.ID != 0 {
				t.Fatal("obj.ID should be 0")
			}
			obj.ID = 1
			obj.Name = "test"

			obj.Free()
			obj.Free()
		}()
	}
	time.Sleep(time.Second)
}
