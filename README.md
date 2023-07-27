## UPDATE 2023-07-23
Note that I'm not updating this repo to keep track with the latest fzf changes, nor dio I apply security updates (some of the dependabot updates seem to be high-risk). If you use this code, you're responsible for doing this yourself!
## End update


This repo is an attempt to create a workable go-library from the amazing [fzf][1] package
from Junegunn Choi.

For more information on the why and how, see [this blogpost][2].

### Performance
The goal of this library is to match the performance of the original fzf.
Where possible, performance enhancing cache has been maintained (and moved into non-global locations, so that multiple fzf objects can live simultaneously).
One of the places where a performance-changing change has been made is in returning the exacts characters that match the hit.
In the original fzf, only the matching lines are returned, and only when displaying these in the console the exact characters are retrieved.
This is obviously faster if one has a complex match that returns thousands of results where we only display a couple in the terminal.
For now fzf-lib always returns machting character positions for all matches.
Benchmarks shows that returning the positions has a 5-10% speed cost; considering the extreme speed of the fzf code, I hardly think anyone will miss these 5-10 percents.
If this turns out to be a bottleneck, in the future we may consider alternative options.

The testing suite contains a benchmark, that fuzzy searches a string in a (repeating) list of quotes of different lengths. On my Macbook Pro M1, I got the following timings (this is from the moment a `Search("hello world")` command is given, until the full result is returned):

| Number of lines to search in | time until result (ms) |
|----------|------|
|1,024|0.709|
|2,048|1.140|
|4,096|2.483|
|8,192|4.787|
|16,384|6.814|
|32,768|12.86|
|65,536|24.92|
|131,072|48.95|
|262,144|95.65|
|524,288|190.9|
|1,048,576|380.1|
|2,097,152|767.5 (0.7s)|
|4,194,304|1,577 (1.5s)|
|8,388,608|3,173 (3s)|
|16,777,216|6,588 (6.5s)|
|33,554,432|33,098 (33s)|

It is hard to properly compare this to `fzf --filter`, since it's not clear how much is overhead / parsing the data, etc.
However the results are in the same order; you can see more info on [my blog][2].
Obviously the results depend a lot on system and exact use, however I do think it's fair to say that up until 100k strings to search in, you should have a performance that will work for updating-the-results-while-typing, especially since that will have caching that was not taken into account here.


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
    myFzf.Search(`^hel owo`)
    result = <- myFzf.GetResultChannel()
    fmt.Printf("%#v", result)
    time.Sleep(200 * time.Millisecond)
    myFzf.Search(`^hy owo`)
    result = <- myFzf.GetResultChannel()
    fmt.Printf("%#v", result)
    myFzf.End() // Note: not strictly necessary since end of program
}
```

Note: If the channel is not being read, the search go routine will block

The following options can be set (most are 1-on-1 matches to fzf commandline optioens with the same name
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
Please see [this blog post][2] for more information.

#### Why does this look nothing like a proper golang module/library
Quite honestly, because this is my first Go project that I'm working on.
I appreciate any feedback on how to make it better.


#### What version of fzf is this code based on
The code was cloned from the fzf master branch on 8 June 2021, latest commit being `7191ebb615f5d6ebbf51d598d8ec853a65e2274d`.
This means that it's basically version `0.27.2`, with some bug fixes.
Fzf follows its own version numbering scheme, where we creep up on v1.0 for a production ready release.
The old git tags (with fzf releases) have been removed from the github repository.

#### What is still missing
The wishlist for v1.0 is (in addition to extra (stress)tests):

- Send SearchProgess messages on the result channel if the search takes more than 200ms, so that a progress bar can be shown
- See if we can automatically call `myFzf.End()` when the item goes out of scope.
- Allow selection of algorithm v1, in case someone would want that.
- Probably some work to make this act nicely in the Go ecosystem.

#### Appreciation / Thank You's / Coffee / Beer
If you like the project, I always appeciate feedback (on [my blog][2]), in the issues or by starring this repository.
I still have a dream that one day I'll be able to have fans pay for my coffees and beers; when that time comes, you'll find a button here!

#### Can you share more info into the inner workings of fzf(-lib)
Fzf-lib is just the code that remains if you strip all interfaces from [fzf][1].
You should approach fzf's author for more info on choices made and how stuff actually works.

Having said that, I was required to do some reverse-engineering of some parts; I documented some of this in response to a github issue: https://github.com/reinhrst/fzf-lib/issues/1.
If you're interested in this kind of stuff, I suggest you read there.

[1]: https://github.com/junegunn/fzf
[2]: https://blog.claude.nl/tech/making-fzf-into-a-golang-library-fzf-lib/
