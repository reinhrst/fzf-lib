package main

import (
	"github.com/junegunn/fzf/src"
    "syscall/js"
    "encoding/json"
)

var version string = "0.27"
var revision string = "devel"

var fzfs []*fzf.Fzf

func NewFzfJs(this js.Value, args []js.Value) interface{} {
    var hayStack [][]byte
    for _, straw := range args {
        hayStack = append(hayStack, []byte(straw.String()))
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
    jsonString, _ := json.Marshal(results)
    return js.ValueOf(string(jsonString))
}

func main() {
    // fzf.ClaudeMatchOnce()
	//fzf.Run(fzf.ParseOptions(), version, revision)
    c := make(chan struct{}, 0)
    js.Global().Set("NewFzf", js.FuncOf(NewFzfJs))
    js.Global().Set("FzfFind", js.FuncOf(Find))
    <-c
}
