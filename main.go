package main

import (
    "github.com/reinhrst/fzf-lib/src"
    "fmt"
)

var version string = "0.27"
var revision string = "devel"

func test() {
    var hayStack [][]byte = [][]byte{
        []byte("apple"),
        []byte("pear"),
        []byte("grape"),
        []byte("apple pear"),
    }
    myFzf := fzf.NewFzf(hayStack)
    results := myFzf.Find([]rune("pe a"))
    for _, result := range results {
        fmt.Printf("%s %d %d %v\n",
            result.Key, result.HayIndex, result.Score, *result.Positions)
    }
}


func main() {
    test()
    //startWasmServer()
}
