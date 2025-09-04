package testdata

import (
	"os"
	"path/filepath"
	"runtime"
)

var TestBookBytes []byte

func init() {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	path := filepath.Join(dir, "test_quotes.epub")

	var err error
	TestBookBytes, err = os.ReadFile(path)
	if err != nil {
		panic(err)
	}
}
