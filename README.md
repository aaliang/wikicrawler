The purpose of this project is to download wiki entries (specifically wikivoyage/wikitravel pages) and store them on disk. Since I find myself without internet access a lot (on the road/abroad) this is particularly useful for myself if I want to read up on something without internet access. (obviously the point is to download these pages first)

- this is not yet finished, only the crawler works at this point

dependencies:

golang
goquery (https://github.com/PuerkitoBio/goquery) - which is awesome by the way
 - $ go get github.com/PuerkitoBio/goquery
 - 

usage:

go run wikicrawler.go [url] [depth]

where [url] is the base page to crawl from
      [depth] is the number of times 'recurse' on from the base page

TODO: - fix goroutine issue, not sure how it will scale for very large depths
      - add additional options for which links to follow
      - add pagewriter class
      - add abstractions, right now it is super procedural
      - filter out querystrings
