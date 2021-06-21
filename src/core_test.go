package fzf

import (
    "testing"
    "sync"
    "reflect"
)

var hayStack [][]byte = [][]byte{
    []byte("apple"),
    []byte("pear"),
    []byte("grape"),
    []byte("apple pear"),
}

func searchHayStack(opts Options, needle string) SearchResult {
    myFzf := New(hayStack, opts)
    var result SearchResult
    var wg sync.WaitGroup
    wg.Add(1)
    go func () {
        defer wg.Done()
        result = <- myFzf.GetResultCannel()
    }()
    myFzf.Search([]rune("pe a"))
    wg.Wait()
    myFzf.End()
    return result
}


func TestSearch(t *testing.T) {
    result := searchHayStack(DefaultOptions(), `pe a`)
    if len(result.Matches) != 4 {
        t.Errorf("Expected 4 results, got %d", len(result.Matches))
    }
}

func TestSearchOrder(t *testing.T) {
    tables := []struct {
        sortCriteria []Criterion
        hits []string
    }{
        {[]Criterion{ByScore, ByLength},[]string{"apple pear", "pear", "apple", "grape"}},
        {[]Criterion{}, []string{"apple", "pear", "grape", "apple pear"}},
    }

    for _, table := range tables {
        options := DefaultOptions()
        options.Sort = table.sortCriteria
        result := searchHayStack(options, `pe a`)
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
