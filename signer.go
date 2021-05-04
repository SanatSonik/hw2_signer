package main

import (
	"fmt"
	"sort"
	"strings"
	"sync"
)

// сюда писать код
func ExecutePipeline(jobs ...job) {

	wg := &sync.WaitGroup{}
	in := make(chan interface{})
	out := make(chan interface{})

	for i := 0; i < len(jobs); i++ {
		wg.Add(1)

		go func(i job, in, out chan interface{}) {
			defer wg.Done()
			defer close(out)
			i(in, out)
		}(jobs[i], in, out)

		in = out
		out = make(chan interface{})

	}
	wg.Wait()
}

func SingleHash(in, out chan interface{}) {
	m := &sync.Mutex{}
	wg := &sync.WaitGroup{}

	for i := range in {
		wg.Add(1)
		go SingleHash2(wg, m, out, fmt.Sprintf("%v", i))
	}
	wg.Wait()
}

func SingleHash2(wg *sync.WaitGroup, m *sync.Mutex, out chan interface{}, data string) {
	defer wg.Done()
	ans1 := make(chan string)
	ans2 := make(chan string)

	go func() {
		defer close(ans1)
		ans1 <- DataSignerCrc32(data)
	}()

	go func() {
		defer close(ans2)
		m.Lock()
		mb5 := DataSignerMd5(data)
		m.Unlock()
		ans2 <- DataSignerCrc32(mb5)
	}()

	out <- fmt.Sprintf("%s~%s", <-ans1, <-ans2)
}

func MultiHash(in, out chan interface{}) {
	m := &sync.Mutex{}
	wg := &sync.WaitGroup{}

	for i := range in {
		wg.Add(1)
		go MultiHash2(wg, m, out, fmt.Sprintf("%v", i))
	}
	wg.Wait()

}

func MultiHash2(wg *sync.WaitGroup, m *sync.Mutex, out chan interface{}, data string) {
	defer wg.Done()
	wg2 := sync.WaitGroup{}

	multiAns := [6]string{}

	for th := 0; th <= 5; th++ {
		wg2.Add(1)
		go func(th int) {
			defer wg2.Done()
			hash := DataSignerCrc32(fmt.Sprintf("%d%s", th, data)) // crc32(th+data)
			m.Lock()
			multiAns[th] = hash
			m.Unlock()
		}(th)
	}

	wg2.Wait()
	out <- strings.Join(multiAns[:], "")
}

func CombineResults(in, out chan interface{}) {
	finalAns := make([]string, 0)

	for i := range in {
		finalAns = append(finalAns, fmt.Sprintf("%v", i))
	}

	sort.Strings(finalAns)
	result := strings.Join(finalAns, "_")
	out <- result

}
