package main

import (
	"sort"
	"strconv"
	"sync"
	"strings"
)

func ExecutePipeline(jobs ...job) {
	in := make(chan interface{})
	var out chan interface{}

	wg := &sync.WaitGroup{}

	for _, currentJob := range jobs {
		out = make(chan interface{})
		wg.Add(1)

		go func(in, out chan interface{}, currentJob job) {
			defer wg.Done()
			currentJob(in, out)
			close(out)
		}(in, out, currentJob)

		in = out
	}
	wg.Wait()
}

func SingleHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}
	mu := &sync.Mutex{}

	for data := range in {
		wg.Add(1)
		go func(data interface{}) {
			defer wg.Done()
			strData := strconv.Itoa(data.(int))

			mu.Lock()
			md5Hash := DataSignerMd5(strData)
			mu.Unlock()

			crc32Chan := make(chan string)
			go func() { crc32Chan <- DataSignerCrc32(strData) }()

			crc32md5 := DataSignerCrc32(md5Hash)
			crc32data := <-crc32Chan

			out <- crc32data + "~" + crc32md5
			close(crc32Chan)
		}(data)
	}
	wg.Wait()
}

func MultiHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}
	for data := range in {
		wg.Add(1)
		go func(data interface{}) {
			defer wg.Done()
			mu := &sync.Mutex{}
			var result [6]string
			mhWg := &sync.WaitGroup{}

			for i := 0; i < 6; i++ {
				mhWg.Add(1)
				go func(i int) {
					defer mhWg.Done()
					res := DataSignerCrc32(strconv.Itoa(i) + data.(string))
					mu.Lock()
					result[i] = res
					mu.Unlock()
				}(i)
			}

			mhWg.Wait()
			out <- strings.Join(result[:], "")
		}(data)
	}
	wg.Wait()
}

func CombineResults(in, out chan interface{}) {
	var results []string

	for data := range in {
		results = append(results, data.(string))
	}
	sort.Strings(results)

	out <- joinResults(results, "_")
}

func joinResults(results []string, separator string) string {
	result := ""
	for i, res := range results {
		if i > 0 {
			result += separator
		}
		result += res
	}
	return result
}
