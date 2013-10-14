package main

import (
    "fmt"
    "github.com/PuerkitoBio/goquery"
    // "io/ioutil"
    "net/url"
    "os"
    "strconv"
    "sync"
)

//for now set this to true

func Scrape(target string) (string, map[string]bool) {

    // Load the HTML document
    var doc *goquery.Document
    var e error

    if doc, e = goquery.NewDocument(target); e != nil {
        panic(e.Error())
    }

    linkSet := make(map[string]bool)

    doc.Find("a").Each(func(i int, s *goquery.Selection) {
        link, _ := s.Attr("href")
        linkSet[link] = true
    })

    return doc.Text(), linkSet
}

func Crawl(url string, depth int, hosturi string) {
    m := map[string]bool{url: true}
    var mx sync.Mutex
    var wg sync.WaitGroup
    var c2 func(string, int)
    c2 = func(url string, depth int) {
        defer wg.Done()
        if depth <= 0 {
            return
        }
        body, urls := Scrape(url)

        fmt.Println(body)

        fmt.Println("found: %s", url)
        mx.Lock()
        for u, _ := range urls {
            if !m[u] {
                m[u] = true
                wg.Add(1)
                go c2(u, depth-1)
            }
        }
        mx.Unlock()
    }
    wg.Add(1)
    c2(url, depth)
    wg.Wait()
}

func main() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Println(r)
        }
    }()

    // handle args
    // TODO: arguments should be a little more robust
    if len(os.Args[1:]) != 2 {
        panic("improper usage.\n args[0] = url; args[1] = depth \n ex. go run wikicrawler.go www.google.com 3")
    }

    target := os.Args[1]
    depth, err := strconv.Atoi(os.Args[2])
    if err != nil {
        panic("improper usage.\n depth (args[1]) should be an integer > 0")
    }

    u, err := url.Parse(target)
    if err != nil {
        panic("improper usage.\n url (args[0]) is not accessible")
    }

    hosturi := u.Scheme + "://" + u.Host

    Crawl(target, depth, hosturi)
    // Scrape(target)
    fmt.Println("done")
}
