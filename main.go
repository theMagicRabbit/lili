package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func main() {
	htmlFile, err := os.Open("samples/sample_list.html")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer htmlFile.Close()

	htmlFileBytes, err := io.ReadAll(htmlFile)
	screwYouLinkedInString := html.UnescapeString(string(htmlFileBytes))
	screwYouLinkedInReader := strings.NewReader(screwYouLinkedInString)

	cleanFile, err := os.Create("samples/sample_list_cleaned.html")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	cleanFile.WriteString(screwYouLinkedInString)
	defer cleanFile.Close()
	
	table, tr, div := 0, 0, 0
	z := html.NewTokenizer(screwYouLinkedInReader)
	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			fmt.Println(z.Err().Error())
			return
		case html.TextToken:
			if div > 0 {
				fmt.Println(string(z.Text()))
			}
		case html.StartTagToken, html.EndTagToken:
			tn, hasAttr := z.TagName()
			if  string(tn) == "table" {
				if tt == html.StartTagToken {
					table++
				} else {
					table--
				}
			} else if table > 0 && string(tn) == "tr" {
				if tt == html.StartTagToken {
					tr++
				} else {
					tr--
				}
			} else if tr > 0 && string(tn) == "div" {
				divContainsData := false
				ParseAttributes:
					for hasAttr {
						attrKey, attrVal, moreAttr := z.TagAttr()
						hasAttr = moreAttr
						switch string(attrKey) {
						case "aria-label":
							if string(attrVal) == "Lead Name" {
								divContainsData = true
								break ParseAttributes
							}
						}
					}
				if tt == html.StartTagToken && divContainsData {
					div++
				} else if tt == html.EndTagToken && div > 0{
					div--
				}
			}
		}
	}
}
