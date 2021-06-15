/*
Package fzf implements fzf, a command-line fuzzy finder.

The MIT License (MIT)

Copyright (c) 2013-2021 Junegunn Choi

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package fzf

import (
    "github.com/junegunn/fzf/src/util"
    "github.com/junegunn/fzf/src/algo"
)

type Fzf struct {
    eventBox *util.EventBox
    matcher *Matcher
    chunkList *ChunkList
    slab *util.Slab
}

func NewFzf(hayStack [][]byte) *Fzf {
    var itemIndex int32
    var chunkList = NewChunkList(func(item *Item, data []byte) bool {
            item.text = util.ToChars(data)
            item.text.Index = itemIndex
            itemIndex++
            return true
        })

    for _, hayStraw := range hayStack {
        chunkList.Push(hayStraw)
        println(string(hayStraw))
    }

    eventBox := util.NewEventBox()
    patternBuilder := func(runes []rune) *Pattern {
        return BuildPattern(
            true, algo.FuzzyMatchV2, true, CaseSmart, true, true,
            true, make([]Range, 0), Delimiter{}, runes)
    }
    matcher := NewMatcher(patternBuilder, true, false, eventBox)

    return &Fzf {
        eventBox,
        matcher,
        chunkList,
        util.MakeSlab(slab16Size, slab32Size),
    }
}

type FzfResult struct {
    Key []rune
    HayIndex int32
    Bonus int
    Positions *[]int
}


func (fzf *Fzf) Find(needle []rune) []FzfResult {
    println(string(needle))
    pattern := fzf.matcher.patternBuilder(needle)
    snapshot, _ := fzf.chunkList.Snapshot()
    merger, _ := fzf.matcher.scan(MatchRequest{
        chunks:  snapshot,
        pattern: pattern})

    var results []FzfResult
    for i := 0; i < merger.Length(); i++ {
        item := merger.Get(i).item
        _, bonus, pos := pattern.extendedMatch(item, true, fzf.slab)
        results = append(results, FzfResult{
            item.text.ToRunes(),
            item.text.Index,
            bonus,
            pos,
        })
        println(item.text.Index)
    }
    return results
}

