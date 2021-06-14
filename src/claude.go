package fzf

import (
    "github.com/junegunn/fzf/src/util"
    "github.com/junegunn/fzf/src/algo"
    "fmt"
)


/*
func ClaudeMatcherLoop() {

    opts := DefaultOptions()
    forward := true
    eventBox := util.NewEventBox()
    sort := opts.Sort > 0
    patternBuilder := func(runes []rune) *Pattern {
        return BuildPattern(
            opts.Fuzzy, opts.FuzzyAlgo, opts.Extended, opts.Case, opts.Normalize, forward,
            opts.Filter == nil, opts.Nth, opts.Delimiter, runes)
    }


    ansiProcessor := func(data []byte) (util.Chars, *[]ansiOffset) {
        return util.ToChars(data), nil
    }
    var itemIndex int32
    var chunkList = NewChunkList(func(item *Item, data []byte) bool {
            item.text, item.colors = ansiProcessor(data)
            item.text.Index = itemIndex
            itemIndex++
            return true
        })

    var reader *Reader
    reader = NewReader(func(data []byte) bool {
        return chunkList.Push(data)
    }, eventBox, opts.ReadZero, opts.Filter == nil)
    go reader.ReadSource()

    eventBox.Unwatch(EvtReadNew)
    eventBox.WaitFor(EvtReadFin)

    go output(eventBox)

    matcher := NewMatcher(patternBuilder, sort, opts.Tac, eventBox)
    go matcher.Loop()

    snapshot, _ := chunkList.Snapshot()
    time.Sleep(100 * time.Millisecond)
    fmt.Println("let's do something")
    time.Sleep(100 * time.Millisecond)

    matcher.reqBox.Set(EvtSearchNew, MatchRequest{
        chunks:  snapshot,
        pattern: patternBuilder([]rune("s"))})
    time.Sleep(300 * time.Millisecond)

    matcher.reqBox.Set(EvtSearchNew, MatchRequest{
        chunks:  snapshot,
        pattern: patternBuilder([]rune("sg"))})
    time.Sleep(300 * time.Millisecond)

    matcher.reqBox.Set(EvtSearchNew, MatchRequest{
        chunks:  snapshot,
        pattern: patternBuilder([]rune("sgo"))})
    time.Sleep(300 * time.Millisecond)

    matcher.reqBox.Set(EvtSearchNew, MatchRequest{
        chunks:  snapshot,
        pattern: patternBuilder([]rune("sgo "))})
    time.Sleep(300 * time.Millisecond)

    matcher.reqBox.Set(EvtSearchNew, MatchRequest{
        chunks:  snapshot,
        pattern: patternBuilder([]rune("sgo ^"))})
    time.Sleep(300 * time.Millisecond)

    matcher.reqBox.Set(EvtSearchNew, MatchRequest{
        chunks:  snapshot,
        pattern: patternBuilder([]rune("sgo ^s"))})
    time.Sleep(300 * time.Millisecond)

    fmt.Println("done")
    time.Sleep(100 * time.Millisecond)
}


func output(eventBox *util.EventBox) {
    var ansi = false
    for {
        eventBox.Wait(func(events *util.Events) {
            for _, val := range *events {
                switch val := val.(type) {
                case *Merger:
                    var merger *Merger = val
                    fmt.Println("--- NEW ---")
                    for i := 0; i < merger.Length(); i++ {
                        fmt.Printf("%d: %s\n",
                        merger.Get(i).item.text.Index,
                        merger.Get(i).item.AsString(ansi))
                    }
                }
            }
            events.Clear()
        })
    }
}
*/

func ClaudeMatchOnce() {

    eventBox := util.NewEventBox()
    patternBuilder := func(runes []rune) *Pattern {
        return BuildPattern(
            true, algo.FuzzyMatchV2, true, CaseSmart, true, true,
            true, make([]Range, 0), Delimiter{}, runes)
    }


    var itemIndex int32
    var chunkList = NewChunkList(func(item *Item, data []byte) bool {
            item.text = util.ToChars(data)
            item.text.Index = itemIndex
            itemIndex++
            return true
        })
    matcher := NewMatcher(patternBuilder, true, false, eventBox)

    chunkList.Push([]byte("hello"))
    chunkList.Push([]byte("world"))
    chunkList.Push([]byte("hello world"))
    chunkList.Push([]byte("helloooo"))
    chunkList.Push([]byte("hellwor"))
    chunkList.Push([]byte("hell_worl"))

    pattern := patternBuilder([]rune("loo"))
    matcher.sort = pattern.sortable

    slab := util.MakeSlab(slab16Size, slab32Size)

    snapshot, _ := chunkList.Snapshot()
    merger, _ := matcher.scan(MatchRequest{
        chunks:  snapshot,
        pattern: pattern})
    for i := 0; i < merger.Length(); i++ {
        item := merger.Get(i).item
        offsets, bonus, pos := pattern.extendedMatch(item, true, slab)

        fmt.Printf("%d: %s %v %d %v\n",
            item.text.Index,
            item.AsString(false),
            offsets,
            bonus,
            pos,
        )
    }
}
