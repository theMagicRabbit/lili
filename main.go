package main

import (
	"fmt"
	"os"
	
	"golang.org/x/net/html"
)

func main() {
	htmlFile, err := os.Open("samples/test.html")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	
	htmlNode, err := html.Parse(htmlFile)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Printf("%#vPrintln\n", htmlNode)
	os.Exit(0)
}
