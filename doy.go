package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/html"
)

func fetchPage(url string) (body string, err error) {
	resp, err := http.Get(url)
	defer resp.Body.Close()
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	body = string(bytes)
	if err != nil {
		log.Fatal(err)
	}
	return body, nil
}

func makeUrl(date time.Time) string {
	return fmt.Sprintf("https://www.daysoftheyear.com/days/%04d/%02d/%02d/", date.Year(), date.Month(), date.Day())
}

func main() {
	now := time.Now()
	url := makeUrl(now)
	data := ""
	var err error
	if data, err = fetchPage(url); err != nil {
		log.Fatal(err)
	}
	var top *html.Node
	top, err = html.Parse(strings.NewReader(data))
	if err != nil {
		log.Fatal(err)
	}
	var descender func(*html.Node)
	descender = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "h3" {
			for _, a := range n.Attr {
				if a.Key == "class" && a.Val == "card-title" {
					var s strings.Builder
					html.Render(&s, n.LastChild.FirstChild)
					fmt.Println(s.String())
					return
				}
			}
		} else {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				descender(c)
			}
		}
	}
	descender(top)
}
