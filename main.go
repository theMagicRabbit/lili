package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
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
	ChanMap        map[int]chan bool
	FinishedChans  int
	TotalChans     int
}

func main() {
	liliState := State{
		LeadDataType: LeadDataNone,
		ConfigFileName: "lili/config.toml",
		LeadMap: make(map[string]Lead),
		LeadMapLock: sync.RWMutex{},
		ChanMap: make(map[int]chan bool, 1), 
		FinishedChans: 0,
	}
	err := liliState.ReadConfig()
	if err != nil {
		log.Fatal(err)
	}

	input, ok := liliState.Config.Directories["input"]
	if !ok {
		log.Fatal("Input directory not listed in config")
	}

	userHome, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	inputDir := path.Join(userHome, input)
	inputFileGlob := fmt.Sprintf("%s/*.html", inputDir)

	inputFiles, err := filepath.Glob(inputFileGlob)
	if err != nil {
		log.Fatal(err)
	}

	if len(inputFiles) == 0 {
		log.Fatalf("No files found in: %s\n", inputDir)
	}

	mapCounter := 0
	for _, f := range inputFiles {
		isFinished := make(chan bool, 1)
		liliState.ChanMap[mapCounter] = isFinished
		go liliState.ProcessHTMLFile(f, isFinished)
		mapCounter++
	}

	liliState.TotalChans = len(liliState.ChanMap)
	for liliState.FinishedChans < len(liliState.ChanMap) {
		for key, ch := range liliState.ChanMap {
			select {
			case _ = <-ch:
				delete(liliState.ChanMap, key)
				liliState.FinishedChans++
			default:
			}
		}
	}

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

func (s *State) ProcessHTMLFile(fileName string, isFinished chan bool) {
	reader, err := ReadHTMLFile("samples/sample_list.html")
	if err != nil {
		log.Println(err)
		return
	}

	s.ParseHTMLListPage(reader)
	isFinished <-true
	close(isFinished)
}
