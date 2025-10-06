package main

import (
	"fmt"
	"os"
	
	"golang.org/x/net/html"
)

func main() {
	htmlFile, err := os.Open("samples/sample_list.html")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	
	depth := 0
	z := html.NewTokenizer(htmlFile)
	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			fmt.Println(z.Err().Error())
			return
		case html.TextToken:
			if depth > 0 {
				fmt.Println(string(z.Text()))
			}
		case html.StartTagToken, html.EndTagToken:
			tn, _ := z.TagName()
			if len(tn) == 1 && tn[0] == 'a' {
				if tt == html.StartTagToken {
					depth++
				} else {
					depth--
				}
			}
		}
	}

}
