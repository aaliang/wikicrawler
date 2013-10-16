The purpose of this project is to download wiki entries (specifically wikivoyage/wikitravel pages) and store them on disk. Since I find myself without internet access a lot (on the road/abroad) this is particularly useful for myself if I want to read up on something without internet access. (obviously the point is to download these pages first)

- this is not yet finished. there are a bunch of issues that are still outstanding
- However this successfully crawls pages to a specified depth, and saves the result onto the filesystem

dependencies:

golang
goquery (https://github.com/PuerkitoBio/goquery) - which is awesome by the way
 - $ go get github.com/PuerkitoBio/goquery

 also, make sure your $GOPATH variable is set!

usage:

go run wikicrawler.go [url] [depth]

ex:
$ go run wikicrawler.go http://wikivoyage.org/wiki/Taipei 2


where [url] is the base page to crawl from
      [depth] is the number of times 'recurse' on from the base page
      
note: this will create (many) directories in the filesystem on the path that wikicrawler.go is

TODO: - fix goroutine issue, not sure how it will scale for very large depths
      - add additional options for which links to follow
      - add pagewriter class
      - add abstractions, right now it is super procedural
      - filter out querystrings
      - add option to output to specified path
