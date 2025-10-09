package main

import (
	"encoding/csv"
	"log"
	"os"
	"sync"

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
	LeadDataType   LeadData
	ConfigFileName string
	Config         *Config
	LeadMap        map[string]Lead
	LeadMapLock    sync.RWMutex
}

func main() {
	liliState := State{
		LeadDataType: LeadDataNone,
		ConfigFileName: "lili/config.toml",
		LeadMap: make(map[string]Lead),
		LeadMapLock: sync.RWMutex{},
	}
	err := liliState.ReadConfig()
	if err != nil {
		log.Fatal(err)
	}

	reader, err := ReadHTMLFile("samples/sample_list.html")
	if err != nil {
		log.Fatal(err)
	}

	liliState.ParseHTMLListPage(reader)

	w := csv.NewWriter(os.Stdout)
	for _, val := range liliState.LeadMap {
		if err := w.Write([]string{val.Name, val.Title, val.CompanyName, val.Geography}); err != nil {
			log.Fatal(err)
		}

		w.Flush()

		if err := w.Error(); err != nil {
			log.Fatal(err)
		}
	}
}
