package fzf

import (
    "fmt"
)


func ClaudeMatchOnce() {
    fzf := NewFzf([][]byte{
        []byte("hello"),
        []byte("world"),
        []byte("hello world"),
        []byte("helloooo"),
        []byte("hellwor"),
        []byte("hell_worl"),
    })

    results := fzf.Find([]rune("loo"))
    for _, result := range results {
        fmt.Printf("%d: %s %v %d\n",
            result.HayIndex,
            string(result.Key),
            result.Positions,
            result.Bonus,
        )
    }
}
