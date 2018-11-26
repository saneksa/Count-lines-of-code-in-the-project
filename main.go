package main

import (
	"github.com/olekukonko/tablewriter"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

var IGNORE_FOLDERS = [7]string{".git", ".idea", "node_modules", "_release", ".vscode", "package-lock.json", "README.md"}

func checkFolders(path string) bool {
	for _, folder := range IGNORE_FOLDERS {
		var pattern = regexp.MustCompile(folder)
		if match := pattern.FindAllStringIndex(path, -1); len(match) > 0 {
			return true
		}
	}
	return false
}

func FilePathWalkDir(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil{
			log.Fatal(err)
		}
		if checkFolders(path) {
			return nil
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func fileReader (file string, wg *sync.WaitGroup, countChan chan int) {
	defer (*wg).Done()
	f, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	temp := strings.Split(string(f), "\n")
	countChan <- len(temp) - 1

}


func main() {
	root := "/home/alex/project/myapp"
	sum:=0
	wg := *new(sync.WaitGroup)
	files, err := FilePathWalkDir(root)
	wg.Add(len(files))
	if err != nil {
		log.Fatal(err)
	}
	countChan := make(chan int, len(files))

	for _, file := range files {
		go fileReader(file, &wg, countChan)
		sum+= <-countChan
	}
	wg.Wait()

	data := [][]string {
		{"Strings in project", strconv.Itoa(sum)},
	}
	table := tablewriter.NewWriter(os.Stdout)
	for _, v := range data {
		table.Append(v)
	}
	table.Render()

 }
