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
)

func main() {

	const maxGoroutines = 5

	totalChan := make(chan int)
	defer close(totalChan)

	goRoutineChan := make(chan int, maxGoroutines)
	defer close(goRoutineChan)

	finish := make(chan bool)
	defer close(finish)

	errors := make(chan error)

	go func() {
		filePath := filepath.Join(".", "urls.txt")
		file, err := os.Open(filePath)
		if err != nil {
			errors <- err
		}
		scanner := bufio.NewScanner(file)
		goRoutinesCounter := 0
		for scanner.Scan() {

			for {
				if maxGoroutines <= goRoutinesCounter {
					goRoutinesCounter -= <-goRoutineChan
				} else {
					break
				}
			}
			goRoutinesCounter += 1
			go func(line string) {
				resp, err := http.Get(line)
				if err != nil {
					errors <- err
				}
				defer resp.Body.Close()

				bytes, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					errors <- err
				}
				count := strings.Count(string(bytes), "Go")
				totalChan <- count
				fmt.Println(fmt.Sprintf("Count for %s: %d", line, count))
				goRoutineChan <- 1
			}(scanner.Text())

		}
		if err := scanner.Err(); err != nil {
			errors <- err
		}

		for {
			if 0 < goRoutinesCounter {
				goRoutinesCounter -= <-goRoutineChan
			} else {
				break
			}
		}
		finish <- true
	}()

	func() {
		var total int
		for {
			select {
			case count := <-totalChan:
				total = total + count
			case err := <-errors:
				log.Fatal(err)
				return
			case <-finish:
				fmt.Printf("Total: %v", total)
				return
			default:
			}
		}
	}()
}
