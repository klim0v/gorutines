package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type Result struct {
	Url   string
	Count int
}

func main() {
	var wg sync.WaitGroup
	out := make(chan Result)

	limit := make(chan struct{}, 5)
	done := make(chan struct{})
	go printFromPipelineAndDone(out, done)

	wg.Add(1)
	go func() {
		wg.Wait()
		close(out)
	}()
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		limit <- struct{}{}
		wg.Add(1)
		go sendToPipeline(scanner.Text(), out, limit, &wg)
	}
	wg.Done()

	<-done
}

func printFromPipelineAndDone(out <-chan Result, done chan<- struct{}) {
	defer func() { done <- struct{}{} }()
	var total int
	for res := range out {
		fmt.Println(fmt.Sprintf("Count for %s: %d", res.Url, res.Count))
		total += res.Count
	}
	fmt.Println(fmt.Sprintf("Total: %d", total))
}

func sendToPipeline(line string, out chan<- Result, limit <-chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	defer func() { <-limit }()
	count, err := countQuantity(line)
	if err != nil {
		log.Println(err)
		return
	}
	result := Result{Url: line, Count: count}
	out <- result
}

func countQuantity(url string) (int, error) {
	var netClient = &http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := netClient.Get(url)
	if err != nil {
		return 0, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Println(err)
		}
	}()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	return strings.Count(string(bytes), "Go"), nil
}
