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
	Title       string
	CompanyName string
	Geography   string
}

type LeadData int

const (
	LeadDataNone LeadData = iota
	LeadDataName
	LeadDataTitle
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
	
	table, tr, dataDiv, dataSpan := 0, 0, 0, 0
	z := html.NewTokenizer(screwYouLinkedInReader)

	leadDirectory:= make(map[string]Lead)

	liliState := State{
		LeadDataType: LeadDataNone,
	}

	moreTokens := true
	var lastUsedName string

	for moreTokens {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			moreTokens = false
		case html.TextToken:
			if dataDiv > 0 || dataSpan > 0 {
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
					lastUsedName = text
				case LeadDataCompanyName:
					if lead, ok := leadDirectory[lastUsedName]; ok {
						lead.CompanyName = text
						leadDirectory[lastUsedName] = lead
					} else {
						leadDirectory[lastUsedName] = Lead{
							CompanyName: text,
						}
					}
				case LeadDataGeography:
					if lead, ok := leadDirectory[lastUsedName]; ok {
						lead.Geography = text
						leadDirectory[lastUsedName] = lead
					} else {
						leadDirectory[lastUsedName] = Lead{
							Geography: text,
						}
					}
				case LeadDataTitle:
					if lead, ok := leadDirectory[lastUsedName]; ok {
						lead.Title = text
						leadDirectory[lastUsedName] = lead
					} else {
						leadDirectory[lastUsedName] = Lead{
							Title: text,
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
				divContainsData := false
				ParseDivAttributes:
					for hasAttr {
						attrKey, attrVal, moreAttr := z.TagAttr()
						hasAttr = moreAttr
						switch string(attrKey) {
						case "aria-label":
							if string(attrVal) == "Lead Name" {
								liliState.LeadDataType = LeadDataName
								divContainsData = true
								break ParseDivAttributes
							}
						case "data-anonymize":
							if string(attrVal) == "job-title" {
								liliState.LeadDataType = LeadDataTitle
								divContainsData = true
								break ParseDivAttributes
							}
						}
					}
				if tt == html.StartTagToken && divContainsData {
					dataDiv++
				} else if tt == html.EndTagToken && dataDiv > 0{
					dataDiv--
				}
			} else if tr > 0 && string(tn) == "span" {
				spanContainsData := false
				ParseSpanAttributes:
					for hasAttr {
						attrKey, attrVal, moreAttr := z.TagAttr()
						hasAttr = moreAttr
						switch string(attrKey) {
						case "data-anonymize":
							if string(attrVal) == "company-name" {
								liliState.LeadDataType = LeadDataCompanyName
								spanContainsData = true
								break ParseSpanAttributes
							}
						}
					}
				if tt == html.StartTagToken && spanContainsData {
					dataSpan++
				} else if tt == html.EndTagToken && dataSpan > 0{
					dataSpan--
				}
			}
		}
	}
	for _, val := range leadDirectory {
		fmt.Println(val.Name, val.Title, val.CompanyName)
	}
}

