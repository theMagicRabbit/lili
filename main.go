package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"sync"
	"time"
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

	output, ok := liliState.Config.Directories["output"]
	if !ok {
		log.Fatal("Output directory not listed in config")
	}

	archive, ok := liliState.Config.Directories["archive"]
	if !ok {
		log.Fatal("Archive directory not listed in config")
	}

	userHome, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	inputDir := path.Join(userHome, input)
	inputFileGlob := fmt.Sprintf("%s/*.html", inputDir)

	outputDir := path.Join(userHome, output)

	archiveDir := path.Join(userHome, archive)

	err = os.MkdirAll(outputDir, 0750)
	if err != nil {
		log.Fatal(err)
	}

	err = os.MkdirAll(archiveDir, 0750)
	if err != nil {
		log.Fatal(err)
	}

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
	for liliState.FinishedChans < liliState.TotalChans {
		for key, ch := range liliState.ChanMap {
			select {
			case _ = <-ch:
				delete(liliState.ChanMap, key)
				liliState.FinishedChans++
			default:
			}
		}
	}

	outFile := fmt.Sprintf("%s.csv", time.Now().Format("2006-01-02-150405"))
	outPath := path.Join(outputDir, outFile)
	log.Printf("Writing CSV data to: %s\n",outPath)

	outfile, err := os.Create(outPath)
	if err != nil {
		log.Fatal(err)
	}
	defer outfile.Close()

	// stdout := csv.NewWriter(os.Stdout)
	outfileW := csv.NewWriter(outfile)
	for _, val := range liliState.LeadMap {
		row := []string{val.Name, val.Title, val.CompanyName, val.Geography}
		if err := outfileW.Write(row); err != nil {
			log.Fatal(err)
		}

		// stdout.Write(row)
	}

	outfileW.Flush()
	if err := outfileW.Error(); err != nil {
		log.Fatal(err)
	}

	// stdout.Flush()
}

func (s *State) ProcessHTMLFile(fileName string, isFinished chan bool) {
	reader, err := ReadHTMLFile(fileName)
	if err != nil {
		log.Println(err)
		return
	}

	s.ParseHTMLListPage(reader)
	isFinished <-true
	close(isFinished)
}
