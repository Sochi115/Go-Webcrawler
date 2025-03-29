package main

import (
	"log"
	"net/url"
	"sort"
	"sync"
)

type config struct {
	maxPages           int
	pages              map[string]int
	baseUrl            *url.URL
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
}

func (cfg *config) crawlPage(rawCurrentURL string) {
	cfg.maxPages -= 1
	if cfg.maxPages <= 0 {
		return
	}
	cfg.wg.Add(1)
	cfg.concurrencyControl <- struct{}{}
	defer cfg.wg.Done()
	defer func() {
		<-cfg.concurrencyControl
	}()

	if !cfg.isDomainEqual(rawCurrentURL) {
		return
	}

	normalizedUrl, err := normalizeUrl(rawCurrentURL)
	if err != nil {
		log.Fatal(err)
	}

	if cfg.addPageVisit(normalizedUrl) {
		log.Printf("Crawling Page: %v", rawCurrentURL)
		html, _ := getHTML(rawCurrentURL)
		urls, _ := getURLsFromHTML(html, cfg.getRawBaseUrlString())

		for _, url := range urls {
			go cfg.crawlPage(url)
		}
	}
}

func (cfg *config) isDomainEqual(rawCurrentURL string) bool {
	currentUrl, err := url.Parse(rawCurrentURL)
	if err != nil {
		log.Fatal(err)
	}
	baseDomain := cfg.baseUrl.Hostname()
	currentDomain := currentUrl.Hostname()

	return baseDomain == currentDomain
}

func (cfg *config) getRawBaseUrlString() string {
	return cfg.baseUrl.Scheme + "://" + cfg.baseUrl.Host
}

func (cfg *config) addPageVisit(normalizedUrl string) (isFirst bool) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	isFirst = true
	if cfg.pages[normalizedUrl] != 0 {
		isFirst = false
	}
	cfg.pages[normalizedUrl] += 1

	return isFirst
}

func (cfg *config) printReport() {
	log.Println("=============================")
	log.Println("REPORT for " + cfg.getRawBaseUrlString())
	log.Println("=============================")
	keys := make([]string, 0, len(cfg.pages))
	for key := range cfg.pages {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool { return cfg.pages[keys[i]] > cfg.pages[keys[j]] })

	for _, k := range keys {
		log.Printf("Found %v external links for page: %v", cfg.pages[k], k)
	}
}
