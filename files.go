package main

import (
	"errors"
	"fmt"
	"log"
	"os"
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

	err = s.ArchiveHTMLFile(fileName)
	if err != nil {
		var pathError *os.PathError
		if errors.As(err, &pathError) {
			log.Printf("Error while deleting processed file %s: %s\n***Please delete above file or it will be processed again on the next run of lili.\n\n", fileName, pathError.Error())
		} else {
			log.Printf("Error on file %s: %s\nManually archive or delete this file\n", fileName, err)
		}
	}

	s.ParseHTMLListPage(reader)
	isFinished <-true
	close(isFinished)
}

func (s *State) ArchiveHTMLFile(fileToArchive string) error {
	fileName := path.Base(fileToArchive)
	archivePath := path.Join(s.ArchiveDir, fileName)

	archiveFile, err := os.Create(archivePath)
	if err != nil {
		return err
	}
	defer archiveFile.Close()

	data, err := os.ReadFile(fileToArchive)
	if err != nil {
		return err
	}

	_, err = archiveFile.Write(data)
	if err != nil {
		return err
	}

	err = os.Remove(fileToArchive)
	if err != nil {
		return err
	}

	return nil
}
