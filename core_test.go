package fzf

import (
	"reflect"
	"testing"
    "io/ioutil"
    "strings"
)

var hayStack = []string{
	`apple`,
	`pear`,
	`grape`,
	`apple pear`,
}

func searchHayStack(opts Options, needles []string) []SearchResult {
	myFzf := New(hayStack, opts)
	var results []SearchResult
    for _, needle := range needles {
        myFzf.Search(needle)
        results = append(results, <-myFzf.GetResultCannel())
    }
	myFzf.End()
	return results
}

func TestSearch(t *testing.T) {
	result := searchHayStack(DefaultOptions(), []string{`pe a`})[0]
	if len(result.Matches) != 4 {
		t.Errorf("Expected 4 results, got %d", len(result.Matches))
	}
}

func TestSearchOrder(t *testing.T) {
	tables := []struct {
		sortCriteria []Criterion
		hits         []string
	}{
		{[]Criterion{ByScore, ByLength}, []string{`apple pear`, `pear`, `apple`, `grape`}},
		{[]Criterion{}, []string{`apple`, `pear`, `grape`, `apple pear`}},
	}

	for _, table := range tables {
		options := DefaultOptions()
		options.Sort = table.sortCriteria
		result := searchHayStack(options, []string{`pe a`})[0]
		var keys []string
		for _, match := range result.Matches {
			keys = append(keys, match.Key)
		}
		if !reflect.DeepEqual(keys, table.hits) {
			t.Errorf("Results do not match, expected %+v, gotten %+v\n%#v\n%#v\n",
				table.hits, keys, table, result)
		}
	}
}

func TestEmptySearch(t *testing.T) {
    searchstrings := []string{"", " ", " ^ ", " ' ", "  ", "   "}
	results := searchHayStack(DefaultOptions(), searchstrings)
    for i, result := range results {
        if (result.Needle != searchstrings[i]) {
            t.Errorf("Result.Needle is not original searchstring: %#v != %#v",
            result.Needle, searchstrings[i])
        }
        if len(result.Matches) != 4 {
            t.Errorf("Expected 4 results, got %d", len(result.Matches))
        }
    }
}

func TestExactVsNonExactSearch(t *testing.T) {
    searchstrings := []string{`pe`, `'pe`}
	results := searchHayStack(DefaultOptions(), searchstrings)
    matchkeys := []string{}
    for _, match := range results[0].Matches {
        matchkeys = append(matchkeys, match.Key)
    }
    if (len(matchkeys) != 4 ||
        matchkeys[0] != `pear` ||
        matchkeys[1] != `apple pear` ||
        matchkeys[2] != `grape` ||
        matchkeys[3] != `apple`) {
            t.Errorf("Unexpected results, got %#v", matchkeys)
    }
    matchkeys = []string{}
    for _, match := range results[1].Matches {
        matchkeys = append(matchkeys, match.Key)
    }
    if (len(matchkeys) != 3 ||
        matchkeys[0] != `pear` ||
        matchkeys[1] != `apple pear` ||
        matchkeys[2] != `grape`) {
            t.Errorf("Unexpected results for exact match, got %#v", matchkeys)
    }
}

func TestQuotes(t *testing.T) {
    quoteBytes, err := ioutil.ReadFile("testdata/quotes.txt")
    if err != nil {
        panic(err)
    }
    quotes := strings.Split(string(quoteBytes), "\n")
    opts := DefaultOptions()
	myFzf := New(quotes, opts)
    var result SearchResult
    myFzf.Search(`hell`)
    result = <-myFzf.GetResultCannel()
    myFzf.Search(`'hell`)
    result = <-myFzf.GetResultCannel()
    if len(result.Matches) != 1 {
        t.Errorf("Expected 1 result, got %d", len(result.Matches))
    }
	myFzf.End()
}
