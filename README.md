This repo is an attempt to create a workable library from the amazing [fzf][1] package
from Junegunn Choi.
The original package contains a number of parts (and I'm sure my description does not do justice to what it *actually* is):

- A fuzzyfinder (having a list of strings, and a search string, returns those strings that match the search string (according to some fuzzy-find rules), ordered in such a way that the best match comes first)
- A command line interface to this fuzzy-finder, that allows for many customizations
- A ncurses/terminal integration that allows interactive fuzzy-finding, seeing previews, hotkeys, etc
- A VIM plugin, that allows fzf to work within (neo)vim

(and probably it does some more things that I don't even know of).
The author of this repo has been a fervent fzf fan for many years, and thinks it's by far the best fuzzy finder out there.

A limitation of fzf is that it's written as a full program.
Over the years I've found myself a number of times looking for a way to integrate fzf into my own programs and (non-terminal) tools.
Fzf's original author has decided not to create a library version of fzf ([see his reasons][2] -- all the right reasons as far as I'm concerned), and rather suggests [running the fzf executable through a spawned process][3].
In my own programs, I've alternated between using exactly that solution, and baking my own extremely simple version of a fuzzy finder (which never feels as good as the original, and doesn't even get to play in the ballpark next door when it comes to speed).

The goal of this repository is to make a version of the fuzzyfinder part of fzf (stripping out as much as possible of the other stuff) and offer it as a go library.
The reasons to go this (above running fzf in a spawned process) are:

- Even though fzf is blazingly fast, when spawning a new process, piping in a long list of options and get back the results for a certain filter, it does take *some* time. This means it feels less than snappy if one is using this as a live-typing filter.
- The spawned process method unfortunately does not give access to exactly which character positions where matched. The live-typing-searches when running fzf interactively in the terminal show feedback on *exactly which* letters in the strings were matched -- there is no way to create this without this position information.
- There are environments when you cannot run a spawned fzf process -- for instance if you're in a webpage.

It should be considered that the fzf code is super fast, has been thouroughly tested (both by automatic tests, and by millions of users using it all over the globe), and uses a simple and intuitive syntax that has never failed me (and is burned in my musscle memory).
Even though it was tempting to implement the fzf fuzzyfinder part from scratch, I opted in the end (because of the reasons in the last sentence) to take the fzf project and strip out all the parts that are not necessary.

The result is a "bare minimum" library that:
- allows initing a fzf object with a list of string. TODO: allow setting of search options
- allows a "find" call on this object, where a string is provided, which returns an array (sorted by *bonus* descending) of:
    - the string matched
    - the index of the matched string (in the original list of strings)
    - the positions of the letters in the string that were matched
    - the bonus (or score) of the match

For now, any change in the list of strings (or search options) means that a new fzf object should be created.
Doing things this way means that there is no need to worry about race conditions.

Do note that the original fzf author [notes][2] one things that might be needed before fzf can be used as a library is to consider "how to revise the code that was written with the assumption that fzf is a short-lived, one-off process".
I have not encountered adverse issues so far, however feel free to post any issues in the issue tracker.

I should note that this is my first Go module, so I'm sure to be making many noob mistakes; feel free to point them out in the issue tracker :).

```
var hayStack [][]byte = {[]byte(`hello world`), []byte(`hyo world`)}
var myFzf = NewFzf{hayStack}
var results = myFzf.find("^hel owo")
println(len(results), results[0].Key, *results[0].Positions)
```


[1]: https://github.com/junegunn/fzf
[2]: https://github.com/junegunn/fzf/pull/1053#issuecomment-330024275
[3]: https://junegunn.kr/2016/02/using-fzf-in-your-program
