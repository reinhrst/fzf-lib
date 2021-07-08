package fzf

import (
	"fmt"
	"github.com/reinhrst/fzf-lib/algo"
	"github.com/reinhrst/fzf-lib/util"
	"math"      // needed only for benchmark
	"math/rand" // needed only for benchmark
	"strings"   // needed only for benchmark
	"time"      // needed only for benchmark
)

type Options struct {
	// If true, each word (separated by non-escaped spaces) is an independent
	// searchterm. If false, all spaces are literal
	Extended bool
	// if true, default is Fuzzy search (' escapes to make exact search)
	// if false, default is exact search (' escapes to make fuzzy search)
	Fuzzy bool
	// CaseRespect, CaseIgnore or CaseSmart
	// CaseSmart matches case insensitive if the needle is all lowercase, else case sensitive
	CaseMode Case
	// set to False to get fzf --literal behaviour:
	// "Do not normalize latin script letters for matching."
	Normalize bool
	// Array with options from {ByScore, ByLength, ByBegin, ByEnd}.
	// Metches will first be sorted by the first element, ties will be sorted by
	// second element, etc.
	// ByScore: Each match is scored (see algo file for more info), higher score
	// comes first
	// ByLength: Shorter match wins
	// ByBegin: Match closer to begin of string wins
	// ByEnd: Match closer to end of string wins
	//
	// If all methods give equal score (including when the Sort slice is empty),
	// the result is sorted by HayIndex, the order in which they appeared in
	// the input.
	Sort []Criterion
}

func DefaultOptions() Options {
	return Options{
		Extended:  true,
		Fuzzy:     true,
		CaseMode:  CaseSmart,
		Normalize: true,
		Sort:      []Criterion{ByScore, ByLength},
	}
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
	eventBox      *util.EventBox
	matcher       *Matcher
	chunkList     *ChunkList
	slab          *util.Slab
	resultChannel chan SearchResult
}

// Creates a new Fzf object, with the given haystack and the given options
func New(hayStack []string, opts Options) *Fzf {
	var itemIndex int32
	var chunkList = NewChunkList(func(item *Item, data []byte) bool {
		item.text = util.ToChars(data)
		item.text.Index = itemIndex
		itemIndex++
		return true
	})

	for _, hayStraw := range hayStack {
		chunkList.Push([]byte(hayStraw))
	}

	eventBox := util.NewEventBox()
	forward := true
	for _, cri := range opts.Sort {
		if cri == ByEnd {
			forward = false
			break
		}
		if cri == ByBegin {
			break
		}
	}
	patternCache := make(map[string]*Pattern)
	patternBuilder := func(needle string) *Pattern {
		return BuildPattern(
			opts.Fuzzy, algo.FuzzyMatchV2, opts.Extended,
			opts.CaseMode, opts.Normalize, forward, needle, opts.Sort,
			&patternCache)
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
		var progress = false
		quit := false
		fzf.eventBox.Wait(func(events *util.Events) {
			for evt, val := range *events {
				switch evt {
				case EvtSearchFin:
					merger = val.(*Merger)
				case EvtSearchProgress:
					// log.Println("search progress, ignoring for now")
					progress = true
				case EvtQuit:
					quit = true
				default:
					panic(fmt.Sprintf("Unexpected type: %T", val))
				}
			}
			events.Clear()
		})
		if progress && merger == nil {
			continue // TODO do something useful here
		}
		if quit {
			break
		}

		var matchResults []MatchResult
		for i := 0; i < merger.Length(); i++ {
			result := merger.Get(i)
			item := result.item
			pos := result.positions
			score := result.score
			matchResults = append(matchResults, MatchResult{
				Key:       item.text.ToString(),
				HayIndex:  item.Index(),
				Score:     score,
				Positions: *pos,
			})
		}

		result := SearchResult{
			Needle:  merger.pattern.originalText,
			Matches: matchResults,
		}
		fzf.resultChannel <- result
	}
}

func (fzf *Fzf) Search(needle string) {
	snapshot, _ := fzf.chunkList.Snapshot()
	fzf.matcher.Reset(snapshot, needle, false, false, true, false)
}

func (fzf *Fzf) End() {
	fzf.matcher.reqBox.Set(EvtQuit, nil)
	fzf.eventBox.Set(EvtQuit, nil)
	close(fzf.resultChannel)
}

// A function that creates a fully internal benchmark, that can be compared in different environments.
// Note that the benchmark is super basic; use the benchmarks in the test suite
// if you can
func RunBasicBenchmark() {
	fruits := []string{`Abiu`, `Açaí`, `Acerola`, `Ackee`, `African cucumber`, `Apple`, `Apricot`, `Avocado`, `Banana`, `Bilberry`, `Blackberry`, `Blackcurrant`, `Black sapote`, `Blueberry`, `Boysenberry`, `Breadfruit`, `Buddha's hand (fingered citron)`, `Cactus pear`, `Canistel`, `Cempedak`, `Cherimoya (Custard Apple)`, `Cherry`, `Chico fruit`, `Cloudberry`, `Coco De Mer`, `Coconut`, `Crab apple`, `Cranberry`, `Currant`, `Damson`, `Date`, `Dragonfruit (or Pitaya)`, `Durian`, `Egg Fruit`, `Elderberry`, `Feijoa`, `Fig`, `Finger Lime (or Caviar Lime)`, `Goji berry`, `Gooseberry`, `Grape`, `Raisin`, `Grapefruit`, `Grewia asiatica (phalsa or falsa)`, `Guava`, `Hala Fruit`, `Honeyberry`, `Huckleberry`, `Jabuticaba`, `Jackfruit`, `Jambul`, `Japanese plum`, `Jostaberry`, `Jujube`, `Juniper berry`, `Kaffir Lime`, `Kiwano (horned melon)`, `Kiwifruit`, `Kumquat`, `Lemon`, `Lime`, `Loganberry`, `Longan`, `Loquat`, `Lulo`, `Lychee`, `Magellan Barberry`, `Mamey Apple`, `Mamey Sapote`, `Mango`, `Mangosteen`, `Marionberry`, `Melon`, `Cantaloupe`, `Galia melon`, `Honeydew`, `Mouse melon`, `Musk melon`, `Watermelon`, `Miracle fruit`, `Monstera deliciosa`, `Mulberry`, `Nance`, `Nectarine`, `Orange`, `Blood orange`, `Clementine`, `Mandarine`, `Tangerine`, `Papaya`, `Passionfruit`, `Peach`, `Pear`, `Persimmon`, `Plantain`, `Plum`, `Prune (dried plum)`, `Pineapple`, `Pineberry`, `Plumcot (or Pluot)`, `Pomegranate`, `Pomelo`, `Purple mangosteen`, `Quince`, `Raspberry`, `Salmonberry`, `Rambutan (or Mamin Chino)`, `Redcurrant`, `Rose apple`, `Salal berry`, `Salak`, `Satsuma`, `Shine Muscat or Vitis Vinifera`, `Sloe or Hawthorn Berry`, `Soursop`, `Star apple`, `Star fruit`, `Strawberry`, `Surinam cherry`, `Tamarillo`, `Tamarind`, `Tangelo`, `Tayberry`, `Tomato`, `Ugli fruit`, `White currant`, `White sapote`, `Yuzu`}
	for i := 10; i < 20; i++ {
		randomizer := rand.New(rand.NewSource(12345))
		numberlines := int(math.Pow(2, float64(i)))
		sentences := produceSentences(fruits, numberlines, randomizer)
		myFzf := New(sentences, DefaultOptions())
        for _, needle := range []string{`hello world`, `green beans`, `i 'll be back`} {
            starttime := time.Now()
            myFzf.Search(needle)
            result, _ := <-myFzf.GetResultCannel()
            endtime := time.Now()
            fmt.Printf("Run of %d sentences took %.3f ms (%d results)\n",
                len(sentences), float64(endtime.Sub(starttime).Microseconds())/1000,
                len(result.Matches))
        }
        myFzf.End()
	}
}

func produceSentences(fruits []string, number int, randomizer *rand.Rand) []string {
	var result []string
	for len(result) < number {
		nrwords := 3 + randomizer.Intn(10)
		var line []string
		for len(line) < nrwords {
			line = append(line, fruits[rand.Intn(len(fruits))])
		}
		result = append(result, strings.Join(line, " "))
	}
	return result
}
