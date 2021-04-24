package main

import (
	"io/ioutil"
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

func countWords(words []string, m map[string]int, wg *sync.WaitGroup) {
	defer wg.Done()

	for _, word := range words {
		m[word]++
	}
}

type key_val struct {
	k string
	v int
}

type ByVal []key_val

func (a ByVal) Len() int           { return len(a) }
func (a ByVal) Less(i, j int) bool { return a[i].v > a[j].v || (a[i].v == a[j].v && a[i].k < a[j].k) }
func (a ByVal) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func main() {
	// read input from input.txt
	d, err := ioutil.ReadFile("input.txt")
	// check errors
	check(err)
	// convert to lower case
	data := strings.ToLower(string(d))
	// split string
	words := strings.Fields(data)

	n := len(words)

	m := make(map[string]int)

	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg.Add(1)
		countWords(words[(n*i)/5:((i+1)*n)/5], m, &wg)
	}

	wg.Wait()

	f, err1 := os.Create("out.txt")

	check(err1)

	var arr []key_val

	for k, v := range m {
		arr = append(arr, key_val{k, v})
	}

	sort.Sort(ByVal(arr))

	for _, a := range arr {
		f.WriteString(a.k + " : " + strconv.Itoa(a.v) + "\n")
	}

}
