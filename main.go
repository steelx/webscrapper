package main

import (
	"crypto/tls"
	"fmt"
	"github.com/steelx/extractlinks"
	"net/http"
	"net/url"
	"os"
)

var (
	urlQueue  = make(chan string)
	config    = &tls.Config{InsecureSkipVerify: true}
	transport = &http.Transport{
		TLSClientConfig: config,
	}
	hasCrawled = make(map[string]bool)
	netClient  *http.Client
)

func init() {
	netClient = &http.Client{
		Transport: transport,
	}
}

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Println("URL is missing, e.g. webscrapper http://js.org/")
		os.Exit(1)
	}

	baseUrl := args[0]

	go func() {
		urlQueue <- baseUrl
	}()

	for href := range urlQueue {
		if isSameDomain(href, baseUrl) && !hasCrawled[href] {
			crawlLink(href, baseUrl)
		}
	}

}

func crawlLink(href, baseUrl string) {
	hasCrawled[href] = true
	fmt.Println("Crawling.. ", href)
	resp, err := netClient.Get(href)
	checkErr(err)
	defer resp.Body.Close()

	links, err := extractlinks.All(resp.Body)
	checkErr(err)

	for _, l := range links {
		absoluteUrl := toFixedUrl(l.Href, baseUrl)
		go func() {
			urlQueue <- absoluteUrl
		}()
	}
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}

func toFixedUrl(href, base string) string {
	uri, err := url.Parse(href)
	if err != nil {
		return ""
	}
	baseUrl, err := url.Parse(base)
	if err != nil {
		return ""
	}
	uri = baseUrl.ResolveReference(uri)
	return uri.String()
}

func isSameDomain(href, baseUrl string) bool {
	uri, err := url.Parse(href)
	if err != nil {
		return false
	}
	parentUri, err := url.Parse(baseUrl)
	if err != nil {
		return false
	}

	if uri.Host != parentUri.Host {
		return false
	}

	return true
}