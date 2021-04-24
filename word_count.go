package main

import (
	"bufio"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func countWords(words []string, ch chan map[string]int) {
	m := make(map[string]int)
	for _, word := range words {
		m[word]++
	}
	ch <- m
	close(ch)
}

type shared_map struct {
	m  map[string]int
	mu sync.Mutex
}

type key_val struct {
	k string
	v int
}

type ByVal []key_val

func (a ByVal) Len() int           { return len(a) }
func (a ByVal) Less(i, j int) bool { return a[i].v > a[j].v || (a[i].v == a[j].v && a[i].k < a[j].k) }
func (a ByVal) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func (sharedMap *shared_map) add(k string, v int) {
	sharedMap.mu.Lock()
	sharedMap.m[k] += v
	sharedMap.mu.Unlock()
}

func reducer(sharedMap *shared_map, ch [5]chan map[string]int) {
	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(in chan map[string]int) {
			defer wg.Done()
			for m := range in {
				for k, v := range m {
					sharedMap.add(k, v)
				}
			}

		}(ch[i])
	}
	wg.Wait()

	f, err := os.Create("out.txt")

	check(err)

	var arr []key_val

	for k, v := range sharedMap.m {
		arr = append(arr, key_val{k, v})
	}

	sort.Sort(ByVal(arr))

	for _, a := range arr {
		f.WriteString(a.k + " : " + strconv.Itoa(a.v) + " \n")
	}

	f.Close()
}

func main() {
	f, err := os.Open("input.txt")

	check(err)

	var words []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		data := strings.ToLower(scanner.Text())
		words = append(words, strings.Split(data, " ")...)
	}

	n := len(words)

	var ch [5]chan map[string]int

	for i := 0; i < 5; i++ {
		ch[i] = make(chan map[string]int)
		go countWords(words[(n*i)/5:((i+1)*n)/5], ch[i])
	}

	sharedMap := shared_map{m: make(map[string]int)}

	reducer(&sharedMap, ch)

}
