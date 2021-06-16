package fzf

import (
    "github.com/reinhrst/fzf-lib/src/util"
    "github.com/reinhrst/fzf-lib/src/algo"
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
    Key string
    HayIndex int32
    Score int
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
        pos := merger.Get(i).positions
        score := merger.Get(i).score
        results = append(results, FzfResult{
            Key: item.text.ToString(),
            HayIndex: item.Index(),
            Score: score,
            Positions: pos,
        })
        println(item.text.Index)
    }
    return results
}

