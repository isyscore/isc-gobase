package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/isyscore/isc-gobase/file"
)

func TestFile(t *testing.T) {
	// file.WriteFile("./sample.txt", "aaa")
	dir, _ := os.Getwd()
	path := filepath.Join(dir, "sample.txt")

	file.AppendFile(path, "ccc")
}

func TestExtract(t *testing.T) {
	dir, _ := os.Getwd()
	path := filepath.Join(dir, "sample.txt")
	p0 := file.ExtractFilePath(path)
	t.Logf("p0: %s", p0)
	n0 := file.ExtractFileName(path)
	t.Logf("n0: %s", n0)
	e0 := file.ExtractFileExt(path)
	t.Logf("e0: %s", e0)
	c0 := file.ChangeFileExt(path, "xyz")
	t.Logf("c0: %s", c0)
}
