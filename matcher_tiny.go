// +build tinygo

package fzf

import (
)

func numCPU() int {
    // tinygo doesn't support getting the number of CPUs, so we just hardcode
    // something reasonable.
    // A strong argument can be made that tinygo will generally not run on
    // state of the art hardware, so 2 seems a reasonable amount.
    // Note that in fzf creates 8 times the number of CPUs as partitions,
    // so a 4 core system will still be able to saturate more than 2 cores.
    return 2
}
