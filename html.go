package main

import (
	"io"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func ReadHTMLFile(fileName string) (*strings.Reader, error) {
	htmlFile, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer htmlFile.Close()

	htmlFileBytes, err := io.ReadAll(htmlFile)
	screwYouLinkedInString := html.UnescapeString(string(htmlFileBytes))
	screwYouLinkedInReader := strings.NewReader(screwYouLinkedInString)

	return screwYouLinkedInReader, nil
}


func (s *State) ParseHTMLListPage(stringReader *strings.Reader) {
	z := html.NewTokenizer(stringReader)

	moreTokens := true
	table, tr, dataDiv, dataSpan, dataTD := 0, 0, 0, 0, 0
	var lastUsedName string
	threadState := LeadDataNone

	for moreTokens {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			moreTokens = false
		case html.TextToken:
			if dataDiv > 0 || dataSpan > 0 || dataTD > 0 {
				text := strings.ReplaceAll(strings.TrimSpace(string(z.Text())), ",", "")
				if text == "" {
					break
				}
				switch threadState {
				case LeadDataName:
					s.LeadMapLock.RLock()
					lead, ok := s.LeadMap[text]
					s.LeadMapLock.RUnlock()
					if ok {
						lead.Name = text
					} else {
						lead = Lead{
							Name: text,
						}
					}
					s.LeadMapLock.Lock()
					s.LeadMap[text] = lead
					s.LeadMapLock.Unlock()
					lastUsedName = text
				case LeadDataCompanyName:
					s.LeadMapLock.RLock()
					lead, ok := s.LeadMap[lastUsedName];
					s.LeadMapLock.RUnlock()
					if ok {
						lead.CompanyName = text
					} else {
						lead = Lead{
							CompanyName: text,
						}
					}
					s.LeadMapLock.Lock()
					s.LeadMap[lastUsedName] = lead
					s.LeadMapLock.Unlock()
				case LeadDataGeography:
					s.LeadMapLock.RLock()
					lead, ok := s.LeadMap[lastUsedName];
					s.LeadMapLock.RUnlock()
					if ok {
						lead.Geography = text
					} else {
						lead = Lead{
							Geography: text,
						}
					}
					s.LeadMapLock.Lock()
					s.LeadMap[lastUsedName] = lead
					s.LeadMapLock.Unlock()
				case LeadDataTitle:
					s.LeadMapLock.RLock()
					lead, ok := s.LeadMap[lastUsedName];
					s.LeadMapLock.RUnlock()
					if ok {
						lead.Title = text
					} else {
						lead = Lead{
							Title: text,
						}
					}
					s.LeadMapLock.Lock()
					s.LeadMap[lastUsedName] = lead
					s.LeadMapLock.Unlock()
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
								threadState = LeadDataName
								divContainsData = true
								break ParseDivAttributes
							}
						case "data-anonymize":
							if string(attrVal) == "job-title" {
								threadState = LeadDataTitle
								divContainsData = true
								break ParseDivAttributes
							}
						}
					}
				if tt == html.StartTagToken && divContainsData {
					dataDiv++
				} else if tt == html.EndTagToken && dataDiv > 0 {
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
								threadState = LeadDataCompanyName
								spanContainsData = true
								break ParseSpanAttributes
							}
						}
					}
				if tt == html.StartTagToken && spanContainsData {
					dataSpan++
				} else if tt == html.EndTagToken && dataSpan > 0 {
					dataSpan--
				}
			} else if tr > 0 && string(tn) == "td" {
				tdContainsData := false
				ParseTDAttributes:
					for hasAttr {
						attrKey, attrVal, moreAttr := z.TagAttr()
						hasAttr = moreAttr
						switch string(attrKey) {
						case "data-anonymize":
							if string(attrVal) == "location" {
								threadState = LeadDataGeography
								tdContainsData = true
								break ParseTDAttributes
							}
						}
					}
				if tt == html.StartTagToken && tdContainsData {
					dataTD++
				} else if tt == html.EndTagToken && dataTD > 0 {
					dataTD--
				}
			}
		}
	}
}

