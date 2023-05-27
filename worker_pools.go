package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

func log(tag string, format string, iface ...interface{}) {
	_format := "%s [%9s] " + format + "\n"

	iface = append(iface, 0, 0)
	copy(iface[2:], iface[:len(iface)-2])
	iface[0] = time.Now().Format("2006-01-02 15:04:05.000")
	iface[1] = tag

	fmt.Printf(_format, iface...)
}

func worker(id int, jobs <-chan int, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	tag := "worker" + strconv.Itoa(id)
	log(tag, "started")

	for j := range jobs {
		log(tag, "picked job: %d", j)
		time.Sleep(time.Second)
		log(tag, "completed job: %d", j)
		results <- j
	}

	log(tag, "terminating")
}

func main() {
	tag := "main"
	log(tag, "program started")

	var wg sync.WaitGroup

	const numJobs = 3
	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)

	const numWorkers = 2
	wg.Add(numWorkers)
	log(tag, "starting %d workers", numWorkers)
	for i := 0; i < numWorkers; i++ {
		go worker(i+1, jobs, results, &wg)
	}

	log(tag, "assigning %d jobs", numJobs)
	for j := 0; j < numJobs; j++ {
		jobs <- (100 + j)
	}
	close(jobs)

	log(tag, "waiting for all jobs to finish")
	wg.Wait()
	close(results)
	log(tag, "all jobs finished")

	log(tag, "reading results")
	for r := range results {
		log(tag, "result: %d", r)
	}

	log(tag, "main finished")
}
