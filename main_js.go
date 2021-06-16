// +build js,wasm

package main

import (
    "github.com/reinhrst/fzf-lib/src"
    "syscall/js"
)

var fzfs []*fzf.Fzf

func NewFzfJs(this js.Value, args []js.Value) interface{} {
    if !this.IsUndefined() {
        panic(`Expect "this" to be undefined`)
    }
    if len(args) != 1 {
        panic(`Expect exactly one argument`)
    }
    length := args[0].Length()
    if (length < 1) {
        panic(`Call fzf with at least one word in the hayStack`)
    }
    var hayStack [][]byte
    for i :=0; i < args[0].Length(); i++ {
        hayStack = append(hayStack, []byte(args[0].Index(i).String()))
    }
    myFzf := fzf.NewFzf(hayStack)
    fzfs = append(fzfs, myFzf)
    return js.ValueOf(len(fzfs) - 1)
}


func Find(this js.Value, args []js.Value) interface{} {
    fzfNr := args[0].Int()
    needle := []rune(args[1].String())
    myFzf := fzfs[fzfNr]
    results := myFzf.Find(needle)
    var simpleResults []interface{}
    for _, result := range results {
        var simplePositions []interface{}
        for _, pos :=  range *result.Positions {
            simplePositions = append(simplePositions, pos)
        }
        simpleResults = append(simpleResults, map[string]interface{} {
            "Key": result.Key,
            "HayIndex": result.HayIndex,
            "Score": result.Score,
            "Positions": simplePositions,
        })
    }
    return simpleResults
}


func startWasmServer () {
    c := make(chan struct{}, 0)
    js.Global().Set("NewFzf", js.FuncOf(NewFzfJs))
    js.Global().Set("FzfFind", js.FuncOf(Find))
    <-c
}
