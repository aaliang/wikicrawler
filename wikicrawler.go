package main

import (
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
)

var fmt = log.New(os.Stdout, "", 0)

/*
 * given a directory and a url, writes the content to a local file with
 * the 0600 permissions
 */
func WriteToFile(relativeUrl string, directory string, content string) {
	if strings.TrimSpace(directory) == "" {
		directory = "."
	}

	pathname := strings.Join([]string{directory, relativeUrl}, "")
	// fmt.Println("path: %s", pathname)
	err := ioutil.WriteFile(pathname, []byte(content), 0777)
	if err != nil { //the parent directory doesn't exist on the filesystem
		//get our filename
		//this isn't pretty, but it works for now. TODO: look if there's a better function for doing this
		splits := strings.Split(relativeUrl, "/")
		filename := splits[len(splits)-1]
		newpath := strings.TrimSuffix(pathname, filename)
		os.MkdirAll(newpath, 0777)

		//try writing again
		fmt.Println("writing file <%s>", pathname)
		ioutil.WriteFile(pathname, []byte(content), 0777)
		//todo error handle + output
	}
}

/*
 * loads an html document given a hosturi, returns the text of the document
 * followed by the links in the page (on the same domain)
 */
func Scrape(target string, hosturi string) (string, map[string]bool) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("panic: %s", r)
		}
	}()

	// fmt.Println("going to %s", target)

	// Load the HTML document
	var doc *goquery.Document
	var e error

	if doc, e = goquery.NewDocument(target); e != nil {
		panic(e.Error())
	}

	linkSet := make(map[string]bool)

	//TODO: on the leaf nodes, won't need to call this
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		link, _ := s.Attr("href")

		if strings.HasSuffix(strings.ToUpper(link), ".JPG") {
			return
		}

		//for now, only accept links on the same domain
		if strings.HasPrefix(link, "//") {
			linkSet["http:"+link] = true
		} else if strings.HasPrefix(link, "/") {
			linkSet[hosturi+link] = true
		} else if strings.HasPrefix(link, hosturi) {
			linkSet[link] = true
		}
	})

	doc.Find("img").Each(func(i int, s *goquery.Selection) {
		link, _ := s.Attr("src")
	})

	html, _ := doc.Html()

	return html, linkSet
}

/*
 * Top level function for handling crawling a page to a given depth
 * links are followed by spawning goroutines
 */
func CrawlHandler(url string, depth int, hosturi string) {
	m := map[string]bool{url: true}
	var mx sync.Mutex
	var wg sync.WaitGroup
	var c2 func(string, int)
	c2 = func(url string, depth int) {
		defer wg.Done()
		if depth <= 0 {
			return
		}
		body, urls := Scrape(url, hosturi)

		//get rid of the http://
		trimmed_url := strings.TrimPrefix(url, "http://")

		//put it in a directory - TODO: we'll want to allow the user to specify the output directory
		WriteToFile(trimmed_url, "./output/", body)

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
	if len(os.Args[1:]) < 2 {
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

	CrawlHandler(target, depth, hosturi)
	// Scrape(target)
	fmt.Println("done")
}
