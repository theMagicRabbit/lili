package main

import (
	"fmt"
	"log"
	"path"
)

func (s *State) DirectorySetup() (input string, output string, archive string, err error) {
	input, ok := s.Config.Directories["input"]
	if !ok {
		return "", "", "", fmt.Errorf("Input directory not listed in config")
	}

	output, ok = s.Config.Directories["output"]
	if !ok {
		return "", "", "", fmt.Errorf("Output directory not listed in config")
	}

	archive, ok = s.Config.Directories["archive"]
	if !ok {
		return "", "", "", fmt.Errorf("Archive directory not listed in config")
	}
	return input, output, archive, nil
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

func (s *State) ArchiveHTMLFile(fileToArchive string) error {
	fileName := path.Base(fileToArchive)

	return nil
}
