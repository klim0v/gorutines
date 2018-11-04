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
	file, err := getFile()
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	wg := new(sync.WaitGroup)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		wg.Add(1)
		go handleUrl(scanner.Text(), wg)
	}
	wg.Wait()
}

func getFile() (file *os.File, err error) {
	filePath := filepath.Join(".", "urls.txt")
	file, err = os.Open(filePath)
	if err != nil {
		return nil, err
	}
	return file, nil
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
