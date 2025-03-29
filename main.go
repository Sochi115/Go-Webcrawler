package main

import (
	"log"
	"net/url"
	"os"
	"strconv"
	"sync"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		log.Fatal("No website provided.")
	} else if len(args) > 4 {
		log.Fatal("Too many arguments provided.")
	} else {
		baseUrl := args[1]
		maxConcurrency, _ := strconv.Atoi(args[2])
		maxPages, _ := strconv.Atoi(args[3])

		rawBaseUrl, _ := url.Parse(baseUrl)

		pages := make(map[string]int)
		concurrencyControl := make(chan struct{}, maxConcurrency)
		crawler := config{baseUrl: rawBaseUrl, pages: pages, concurrencyControl: concurrencyControl, mu: &sync.Mutex{}, wg: &sync.WaitGroup{}, maxPages: maxPages}

		crawler.wg.Add(1)

		crawler.crawlPage(baseUrl)
		crawler.wg.Done()

		crawler.wg.Wait()
		crawler.printReport()
	}
}
