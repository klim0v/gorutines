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
)

type Result struct {
	Url   string
	Count int
}

func main() {
	var wg sync.WaitGroup
	out := make(chan Result)

	limit := make(chan interface{}, 5)
	wg.Add(1)
	go func() {
		wg.Wait()
		close(out)
	}()
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		wg.Add(1)
		go sendToPipeline(scanner.Text(), out, limit, &wg)
	}
	wg.Done()
	var total int
	for res := range out {
		fmt.Println(fmt.Sprintf("Count for %s: %d", res.Url, res.Count))
		total += res.Count
	}
	fmt.Println(fmt.Sprintf("Total: %d", total))
}

func sendToPipeline(line string, out chan<- Result, limit chan interface{}, wg *sync.WaitGroup) {
	defer wg.Done()
	defer func() {
		<-limit
	}()
	limit <- struct{}{}
	out <- countQuantity(line)
}

func countQuantity(url string) (res Result) {
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Println(err)
		}
	}()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	res.Url = url
	res.Count = strings.Count(string(bytes), "Go")
	return
}
