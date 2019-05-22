package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
	xj "github.com/basgys/goxml2json"
)

func main() {
	start := time.Now()
	ch := make(chan string)
	for _, url := range os.Args[1:] {
		go fetch(url, ch)
	}
	for range os.Args[1:] {
		fmt.Println(<-ch)
	}
	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
}

func fetch(url string, ch chan<- string) {
	start := time.Now()
	resp, err := http.Get(url)
	if err !=nil {
		ch <- fmt.Sprint(err)
		return
	} 
	xssRaw, err := xj.Convert(resp.Body)
  	if err != nil {
		ch <- fmt.Sprintf("While converting rss to json %s: %v", url, err)
		return
  	}
	resp.Body.Close()
	if err != nil {
		ch <- fmt.Sprintf("While reading %s: %v", url, err)
		return
	}
	var xss xssFeed;
	if err := json.Unmarshal(xssRaw.Bytes(), &xss); err != nil {
        ch <- fmt.Sprintf("While converting to object %s: %v", url, err)
		return
    }
	secs := time.Since(start).Seconds()
	ch <- fmt.Sprintf("%.2fs %s %s", secs, xss.Rss.Channel.Item[0], url);
}   


type xssFeed struct {
	Rss struct {
		Dc      string `json:"-dc"`
		Content string `json:"-content"`
		Atom    string `json:"-atom"`
		Channel struct {
			Item []struct {
				Description string `json:"description"`
				Encoded     string `json:"encoded"`
				PubDate     string `json:"pubDate"`
				Category    struct {
					Content string `json:"#content"`
					Domain  string `json:"-domain"`
				} `json:"category"`
				Creator string `json:"creator"`
				GUID    struct {
					Content     string `json:"#content"`
					IsPermaLink string `json:"-isPermaLink"`
				} `json:"guid"`
				Title string `json:"title"`
				Link  string `json:"link"`
			} `json:"item"`
			Description   string `json:"description"`
			Language      string `json:"language"`
			LastBuildDate string `json:"lastBuildDate"`
			Image         struct {
				URL   string `json:"url"`
				Title string `json:"title"`
				Link  string `json:"link"`
			} `json:"image"`
			Title     string        `json:"title"`
			Link      []interface{} `json:"link"`
			TTL       string        `json:"ttl"`
			Copyright string        `json:"copyright"`
		} `json:"channel"`
		Version string `json:"-version"`
	} `json:"rss"`
}
