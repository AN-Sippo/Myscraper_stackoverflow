package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	const root_url = "https://ja.stackoverflow.com"
	const search string = "ML"
	const url = "https://ja.stackoverflow.com/search?tab=Relevance&pagesize=50&q=" + search
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	if res.StatusCode != 200 {
		log.Fatalf("status code with %d", res.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	ans := make([]string, 0)
	doc.Find("#mainbar > div.flush-left.js-search-results > div").Children().Each(func(i int, s *goquery.Selection) {
		u := s.Find("a")
		fmt.Println(i)
		uurl, exists := u.Attr("href")
		if !exists {
			log.Fatal("Could not find the url")
		}
		ans = append(ans, root_url+uurl)

	})

	f, err := os.Create("out.csv")
	f.Write([]byte{0xEF, 0xBB, 0xBF})
	ff := csv.NewWriter(f)
	if err != nil {
		log.Fatal(err)
	}
	for iter, url := range ans {
		fmt.Printf("fetching page %d", iter)
		fmt.Println()
		ans := url_to_question(url)
		// fmt.Println(ans)
		var line []string = []string{ans}
		ff.Write(line)

		time.Sleep(1 * time.Second)
	}
	f.Close()

}

func url_to_question(url string) string {
	const selector string = "#question > div.post-layout > div.postcell.post-layout--right > div.s-prose.js-post-body"
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	if res.StatusCode != 200 {
		log.Fatal(res.StatusCode)
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	const ng_class string = "s-code-block"
	var ans string = ""

	doc.Find(selector).Children().Each(func(i int, s *goquery.Selection) {
		// classname, exists := s.Attr("class")
		// fmt.Println(classname)
		// var flag bool = true
		// if exists && len(classname) > len(ng_class) {
		// 	flag = !(classname[len(classname)-len(ng_class):] == ng_class)
		// }
		// if flag {
		// 	ans = ans + s.Text()
		// }
		// fmt.Println(goquery.NodeName(s))
		if goquery.NodeName(s) != "pre" {
			ans = ans + s.Text()
		}
	})
	return ans
}
