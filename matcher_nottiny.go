// +build !tinygo

package fzf

import (
	"runtime"
)

func numCPU() int {
        return runtime.NumCPU()
}
