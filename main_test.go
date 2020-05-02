package main

import (
	"github.com/linbaozhong/model-gen/cmd"
	"os"
	"path/filepath"
	"testing"
)

func TestModel(t *testing.T) {
	path, _ := os.Getwd()
	filename := filepath.Join(path, "user_model.go")
	cmd.HandleFile(filename)
}
