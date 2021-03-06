package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	start := time.Now()
	ch := make(chan string)

	fmt.Printf("Checking sites for: %q\n", os.Args[2])

	sourceList, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not open source list: $v\n", err)
	}

	var totalSources int
	input := bufio.NewScanner(sourceList)
	for input.Scan() {
		totalSources++
		go fetch(input.Text(), ch)
	}
	sourceList.Close()

	fmt.Println("")

	for i := 0; i < totalSources; i++ {
		fmt.Println(<-ch)
	}
	fmt.Printf("\n%.2fs elapsed\n", time.Since(start).Seconds())
}

func fetch(url string, ch chan<- string) {
	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		ch <- fmt.Sprint(err)
		return
	}
	defer resp.Body.Close()

	responseString, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		ch <- fmt.Sprintf("while raeding %s: %v", url, err)
		return
	}
	secs := time.Since(start).Seconds()

	if strings.Contains(string(responseString), os.Args[2]) {
		ch <- fmt.Sprintf("\033[0;33m%.2fs\033[0m %s \t \033[0;32m=> Found\033[0m", secs, url)
	} else {
		ch <- fmt.Sprintf("\033[0;33m%.2fs\033[0m %s", secs, url)
	}

	ch <- fmt.Sprintf("\033[0;33m%.2fs\033[0m %s", secs, url)

}
