package main

import (
	"fmt"
	"sync"
)

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

//visitedUrls is used to cache urls visited. the value of map is expected to be always true
var visitedUrls = struct {
	sync.RWMutex
	urls map[string]bool
}{urls: make(map[string]bool)}

//to wait for all concurrent crawlers finish their works
var wg sync.WaitGroup

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher) {
	// The implementation of Fetch URLs in parallel done
	// The implementation of "Don't fetch the same URL twice." is done
	defer wg.Done()
	visitedUrls.RLock()
	_, visited := visitedUrls.urls[url]
	visitedUrls.RUnlock()

	if visited {
		return
	}
	if depth <= 0 {
		return
	}

	body, urls, err := fetcher.Fetch(url)
	visitedUrls.Lock()
	visitedUrls.urls[url] = true
	visitedUrls.Unlock()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("found: %s %q\n", url, body)
	for _, u := range urls {
		wg.Add(1)
		go Crawl(u, depth-1, fetcher)
	}
	return
}

func main() {
	populateFakeFatcherStruct()
	wg.Add(1)
	Crawl("http://golang.org/", 4, fetcher)
	wg.Wait()
	fmt.Print("crawled all successfully")
}

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcherStruct struct {
	fakeFetcher map[string]*fakeResult
	sync.RWMutex
}

type fakeResult struct {
	body string
	urls []string
}

func (f fakeFetcherStruct) Fetch(url string) (string, []string, error) {

	f.RLock()
	res, ok := f.fakeFetcher[url]
	f.RUnlock()

	if ok {
		return res.body, res.urls, nil
	}

	return "", nil, fmt.Errorf("not found: %s", url)
}

// fetcher is a populated fakeFetcher.
var fetcher = fakeFetcherStruct{}

func populateFakeFatcherStruct() {
	fetcher.fakeFetcher = make(map[string]*fakeResult)

	fetcher.fakeFetcher["http://golang.org/"] = &fakeResult{"The Go Programming Language",
		[]string{
			"http://golang.org/pkg/",
			"http://golang.org/cmd/",
		},
	}

	fetcher.fakeFetcher["http://golang.org/pkg/"] = &fakeResult{
		"Packages",
		[]string{
			"http://golang.org/",
			"http://golang.org/cmd/",
			"http://golang.org/pkg/fmt/",
			"http://golang.org/pkg/os/",
		},
	}

	fetcher.fakeFetcher["http://golang.org/pkg/fmt/"] = &fakeResult{
		"Package fmt",
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	}

	fetcher.fakeFetcher["http://golang.org/pkg/os/"] = &fakeResult{
		"Package os",
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	}
}
