package fzf

import (
	"fmt"
	"github.com/reinhrst/fzf-lib/src/algo"
	"github.com/reinhrst/fzf-lib/src/util"
	"log"
)

type Options struct {
	Sort []Criterion
}

type SearchResult struct {
	Needle        string
	SearchOptions Options
	Matches       []MatchResult
}

type MatchResult struct {
	Key       string
	HayIndex  int32
	Score     int
	Positions []int
}

type Fzf struct {
	eventBox     *util.EventBox
	matcher      *Matcher
	chunkList    *ChunkList
	slab         *util.Slab
	resultChannel chan SearchResult
}

// Creates a new Fzf object, with the given haystack and the given options
func New(hayStack [][]byte, options Options) *Fzf {
	var itemIndex int32
	var chunkList = NewChunkList(func(item *Item, data []byte) bool {
		item.text = util.ToChars(data)
		item.text.Index = itemIndex
		itemIndex++
		return true
	})

	for _, hayStraw := range hayStack {
		chunkList.Push(hayStraw)
	}

	eventBox := util.NewEventBox()
	patternBuilder := func(runes []rune) *Pattern {
		return BuildPattern(
			true, algo.FuzzyMatchV2, true, CaseSmart, true, true,
			true, make([]Range, 0), Delimiter{}, runes)
	}
	matcher := NewMatcher(patternBuilder, true, false, eventBox)
	resultChannel := make(chan SearchResult)

	fzf := &Fzf{
		eventBox,
		matcher,
		chunkList,
		util.MakeSlab(slab16Size, slab32Size),
		resultChannel,
	}
	fzf.start()
	return fzf
}

func (fzf *Fzf) start() {
	go fzf.loop()
	go fzf.matcher.Loop()
}

func (fzf *Fzf) GetResultCannel() <-chan SearchResult {
    return fzf.resultChannel
}

func (fzf *Fzf) loop() {
	for {
		var merger *Merger
		quit := false
		fzf.eventBox.Wait(func(events *util.Events) {
			for evt, val := range *events {
				switch evt {
				case EvtSearchFin:
					merger = val.(*Merger)
				case EvtSearchProgress:
					log.Println("search progress, ignoring for now")
				case EvtQuit:
					quit = true
				default:
					panic(fmt.Sprintf("Unexpected type: %T", val))
				}
			}
			events.Clear()
		})
		if quit {
			log.Println("Quiting fzf loop")
			break
		}

		var matchResults []MatchResult
		for i := 0; i < merger.Length(); i++ {
			item := merger.Get(i).item
			pos := merger.Get(i).positions
			score := merger.Get(i).score
			matchResults = append(matchResults, MatchResult{
				Key:       item.text.ToString(),
				HayIndex:  item.Index(),
				Score:     score,
				Positions: *pos,
			})
		}

		result := SearchResult{
			Needle:  merger.pattern.cacheKey,
			Matches: matchResults,
		}
		select {
		case fzf.resultChannel <- result:
			log.Println("result sent")
		default:
			log.Println("No listener on the channel")
		}
	}
}

func (fzf *Fzf) Search(needle []rune){
	snapshot, _ := fzf.chunkList.Snapshot()
    fzf.matcher.Reset(snapshot, needle, false, false, true, false)
}
