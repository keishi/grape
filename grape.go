package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sync"
)

var root string
var query *regexp.Regexp
var found = 1
var wg sync.WaitGroup

func readFile(wg *sync.WaitGroup, path string) {
	defer wg.Done()

	buf, err := ioutil.ReadFile(path)
	if err != nil {
		log.Println("Failed to read file:", path, err)
		return
	}
	m := query.Find(buf)
	if m != nil {
		fmt.Println(path)
		found = 0
	}
}

func main() {
	flag.Parse()
	root = flag.Arg(1)

	var err error

	query, err = regexp.Compile(flag.Arg(0))
	if err != nil {
		log.Fatalln(err)
	}

	filepath.Walk(root, func(path string, file os.FileInfo, err error) error {
		if !file.IsDir() {
			wg.Add(1)
			go readFile(&wg, path)
		}
		return nil
	})
	wg.Wait()
	defer os.Exit(found)
}
