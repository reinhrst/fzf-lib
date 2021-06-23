This repo is an attempt to create a workable go-library from the amazing [fzf][1] package
from Junegunn Choi.
The original package contains a number of parts (and I'm sure my description does not do justice to what it *actually* is):

- A fuzzyfinder (having a list of strings, and a search string, returns those strings that match the search string (according to some fuzzy-find rules), ordered in such a way that the best match comes first)
- A command line interface to this fuzzy-finder, that allows for many customizations
- A ncurses/terminal integration that allows interactive fuzzy-finding, seeing previews, hotkeys, etc
- A VIM plugin, that allows fzf to work within (neo)vim

The author of this repo has been a fervent fzf fan for many years, and thinks it's by far the best fuzzy finder out there.

A limitation of fzf is that it's written as a full program.
Over the years I've found myself a number of times looking for a way to integrate fzf into my own programs and (non-terminal) tools.
Fzf's original author has decided not to create a library version of fzf ([see his reasons][2] -- all the right reasons as far as I'm concerned), and rather suggests [running the fzf executable through a spawned process][3].
In my own programs, I've alternated between using exactly that solution, and baking my own extremely simple version of a fuzzy finder (which never feels as good as the original, and doesn't even get to play in the ballpark next door when it comes to speed).

The goal of this repository is to make a version of the fuzzyfinder part of fzf (stripping out as much as possible of the other stuff) and offer it as a go library.
The reasons to go this (above running fzf in a spawned process) are:

- Even though fzf is blazingly fast, when spawning a new process, piping in a long list of options and get back the results for a certain filter, it does take *some* time. This means it feels less than snappy if one is using this as a live-typing filter.
- The spawned process method unfortunately does not give access to exactly which character positions where matched. The live-typing-searches when running fzf interactively in the terminal show feedback on *exactly which* letters in the strings were matched -- there is no way to create this without this position information.
- There are environments when you cannot run a spawned fzf process -- for instance if you're in a webpage (yes, this project was sparked by an interest to compile fzf to Web Assembly and use clientside in a webpage)

It should be considered that the fzf code is super fast, has been thouroughly tested (both by automatic tests, and by millions of users using it all over the globe), and uses a simple and intuitive syntax that has never failed me (and is burned in my musscle memory).
Even though it was tempting to implement the fzf fuzzyfinder part from scratch, I opted in the end (because of the reasons in the last sentence) to take the fzf project and strip out all the parts that are not necessary.
The hope is that this library will maintain the speed and accuracy of the original implementation, while working as a full library.
Since quite some code had to be rewriten to get the library to work(the original fzf author [notes][2] that one things that might be needed before fzf can be used as a library is to consider "how to revise the code that was written with the assumption that fzf is a short-lived, one-off process") , only time will tell if I succeeded in this goal :).

The resulting library:
- allows initing a fzf stuct with a list of strings and search options (note that fzf contains very many (commandline) options, most of them are to control other things than the actual search.
- A `Search` method which starts a search. Results are returned through a channel. If a new `Search` is started before the old returns, the old search is cancelled.
- The data on the channel is either a SearchProgress (if the search takes more than 200ms), or a SearchResult
- A SearchResult is a struct with the SearchKey and Options used, and a list of MatchResults:
    - the string matched
    - the index of the matched string (in the original list of strings)
    - the positions of the letters in the string that were matched
    - the bonus (or score) of the match
- An `End` method. Calling this method is necessary so that the internal go routines are stopped and the Fzf object (which may be quite large since it caches all strings) can be removed from memory.

For now, any change in the list of strings (or search options) means that a new fzf object should be created.


I should note that this is my first Go module, so I'm sure to be making many noob mistakes; feel free to point them out in the issue tracker :).
Obviously also bugs and other issues are welcome.

### Performance
The goal of this library is to match the performance of the original fzf.
Where possible, performance enhancing cache has been maintained (and moved into non-global locations, so that multiple fzf objects can live simultaneously).
One of the places where a performance-changing change has been made is in returning the exacts characters that match the hit.
In the original fzf, only the matching lines are returned, and only when displaying these in the console the exact characters are retrieved.
This is obviously faster if one has a complex match that returns thousands of results where we only display a couple in the terminal.
For now fzf-lib always returns machting character positions for all matches.
If this turns out to be a bottleneck, in the future we may consider alternative options.

Another item that might be a problem (although I doubt it)...
The original fzf wants `[][]byte` to build up its words to search in and then `[]rune` for searching. 
This wrapper makes all input and output in `string`.
Converting should not take much (if any) time, but maybe for very very large datasets there may be a small hit


### Installation

TODO

### Usage


```go
package main

import (
    "fmt"
    "github.com/reinhrst/fzf-lib"
    "time"
    "sync"
)

func main() {
    var options = fzf.DefaultOptions()
    // update any options here
    var hayStack = []string{`hello world`, `hyo world`}
    var myFzf = fzf.New(hayStack, options)
    var result fzf.SearchResult
    wg := sync.WaitGroup{}
    wg.Add(1)
    go func() {
        defer wg.Done()
        outputs := 2
        for {
            result = <- myFzf.GetResultCannel()
            fmt.Printf("%#v", result)
            outputs--
            if outputs == 0 {
                break
            }
        }
    }()
    myFzf.Search(`^hel owo`)
    time.Sleep(200 * time.Millisecond)
    myFzf.Search(`^hy owo`)
    wg.Wait()
    myFzf.End() // Note: not necessary since end of program
}
```

Note: If the channel is not being read, searching will block

The following options can be set
```go
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

```
The DefaultOptions are as follows:
```go
func DefaultOptions() Options {
    return Options{
        Extended: true,
        Fuzzy: true,
        CaseMode: CaseSmart,
        Normalize: true,
        Sort: []Criterion{ByScore, ByLength},
    }
}
```


### FAQ

#### Where are all other useful options
The idea is that most other useful options in fzf cli are easy to build in your client program.
Things like previews, keybindings, and _nth_ item for parsing should be handled by the client code.
Just give the string to match to fzf-lib, and keep your own array with the data you want to show/preview/return and use the index of the matched items to retrieve the full item.
The idea was only to include the items that have to do with _searching and matching_ in this library.
If you feel that I missed an important option, feel free to file an issue.


[1]: https://github.com/junegunn/fzf
[2]: https://github.com/junegunn/fzf/pull/1053#issuecomment-330024275
[3]: https://junegunn.kr/2016/02/using-fzf-in-your-program
