package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func main() {
	lines, err := getUrls()
	if err != nil {
		log.Fatal(err)
	}

	wg := new(sync.WaitGroup)
	for _, line := range lines {
		wg.Add(1)
		go handleUrl(line, wg)
	}
	wg.Wait()
}

func handleUrl(line string, waiteGroup *sync.WaitGroup) {
	defer waiteGroup.Done()
	resp, err := http.Get(line)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(fmt.Sprintf("Count for %s: %d", line, strings.Count(string(bytes), "Go")))
	}
}

func getUrls() (lines []string, err error) {
	filePath := filepath.Join(".", "urls.txt")
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, nil
}
