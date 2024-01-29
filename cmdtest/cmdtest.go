package main

import (
	// "context"

	"os"

	// cliutil "github.com/filecoin-project/lotus/cli/util"
	logging "github.com/ipfs/go-log/v2"
)

// const (
// 	ss32GiB = 32 << 30
// )

var log = logging.Logger("recover")

func main() {

}

func mkdirAll(path string) {
	_, err := os.Stat(path)
	notexist := os.IsNotExist(err)

	if notexist {
		_ = os.MkdirAll(path, 0755) //nolint: gosec
	}
}
