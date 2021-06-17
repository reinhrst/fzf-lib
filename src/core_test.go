package fzf

import (
    "testing"
    "sync"
)


func TestSearch(t *testing.T) {
	var hayStack [][]byte = [][]byte{
		[]byte("apple"),
		[]byte("pear"),
		[]byte("grape"),
		[]byte("apple pear"),
	}
    myFzf := New(hayStack, Options{Sort: []Criterion{ByScore}})
    var result SearchResult
    var wg sync.WaitGroup
    wg.Add(1)
    go func () {
        defer wg.Done()
        result = <- myFzf.GetResultCannel()
    }()
    myFzf.Search([]rune("pe a"))
    wg.Wait()
    if len(result.Matches) != 4 {
        t.Errorf("Expected 4 results, got %d", len(result.Matches))
    }
}
