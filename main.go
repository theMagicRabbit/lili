package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/net/html"
)

type Lead struct {
	Name        string
	CompanyName string
	Geography   string
}

type LeadData int

const (
	LeadDataNone LeadData = iota
	LeadDataName
	LeadDataCompanyName
	LeadDataGeography
)


type State struct {
	LeadDataType LeadData
}

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

	leadDirectory:= make(map[string]Lead)

	liliState := State{
		LeadDataType: LeadDataNone,
	}

	moreTokens := true
	for moreTokens {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			moreTokens = false
		case html.TextToken:
			if div > 0 {
				text := strings.TrimSpace(string(z.Text()))
				if text == "" {
					break
				}
				switch liliState.LeadDataType {
				case LeadDataName:
					if lead, ok := leadDirectory[text]; ok {
						lead.Name = text
						leadDirectory[text] = lead
					} else {
						leadDirectory[text] = Lead{
							Name: text,
						}
					}
				case LeadDataCompanyName:
					if lead, ok := leadDirectory[text]; ok {
						lead.CompanyName = text
						leadDirectory[text] = lead
					} else {
						leadDirectory[text] = Lead{
							CompanyName: text,
						}
					}
				case LeadDataGeography:
					if lead, ok := leadDirectory[text]; ok {
						lead.Geography = text
						leadDirectory[text] = lead
					} else {
						leadDirectory[text] = Lead{
							Geography: text,
						}
					}
				}
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
				divContainsLeadName := false
				ParseAttributes:
					for hasAttr {
						attrKey, attrVal, moreAttr := z.TagAttr()
						hasAttr = moreAttr
						switch string(attrKey) {
						case "aria-label":
							if string(attrVal) == "Lead Name" {
								liliState.LeadDataType = LeadDataName
								divContainsLeadName = true
								break ParseAttributes
							}
						}
					}
				if tt == html.StartTagToken && divContainsLeadName {
					div++
				} else if tt == html.EndTagToken && div > 0{
					div--
				}
			}
		}
	}
	for _, val := range leadDirectory {
		fmt.Println(val.Name)
	}
}

