package main

import (
	"crypto/tls"
	"fmt"
	"github.com/steelx/extractlinks"
	"github.com/steelx/webscrapper/graph"
	"net/http"
	"net/url"
	"os"
	"os/signal"
)

var (
	urlQueue  = make(chan string)
	config    = &tls.Config{InsecureSkipVerify: true}
	transport = &http.Transport{
		TLSClientConfig: config,
	}
	hasCrawled = make(map[string]bool)
	netClient  *http.Client
	graphMap   = graph.NewGraph()
)

func init() {
	netClient = &http.Client{
		Transport: transport,
	}
	go SignalHandler(make(chan os.Signal, 1))
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
		if !hasCrawled[href] {
			crawlLink(href)
		}
	}

}

func SignalHandler(c chan os.Signal) {
	signal.Notify(c, os.Interrupt)
	for s := <-c; ; s = <-c {
		switch s {
		case os.Interrupt:
			fmt.Println("^C received")
			fmt.Println("<----------- ----------- ----------- ----------->")
			fmt.Println("<----------- ----------- ----------- ----------->")
			graphMap.CreatePath("https://youtube.com/jsfunc", "https://youtube.com/YouTubeRedOriginals")
			os.Exit(0)
		case os.Kill:
			fmt.Println("SIGKILL received")
			os.Exit(1)
		}
	}
}

func crawlLink(baseHref string) {
	graphMap.AddVertex(baseHref)
	hasCrawled[baseHref] = true
	fmt.Println("Crawling.. ", baseHref)
	resp, err := netClient.Get(baseHref)
	checkErr(err)
	defer resp.Body.Close()

	links, err := extractlinks.All(resp.Body)
	checkErr(err)

	for _, l := range links {
		if l.Href == "" {
			continue
		}
		fixedUrl := toFixedUrl(l.Href, baseHref)
		if baseHref != fixedUrl {
			graphMap.AddEdge(baseHref, fixedUrl)
		}
		go func(url string) {
			urlQueue <- url
		}(fixedUrl)
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
	if err != nil || uri.Scheme == "mailto" || uri.Scheme == "tel" {
		return base
	}
	baseUrl, err := url.Parse(base)
	if err != nil {
		return ""
	}
	uri = baseUrl.ResolveReference(uri)
	return uri.String()
}
