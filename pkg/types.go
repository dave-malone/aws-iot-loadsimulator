package awsiotloadsimulator

import (
	"fmt"
	"sync"
	"time"
)

type workerFunc func(int) error

type SimulationRequest struct {
	StartClientNumber int `json:"start_client_num"`
	ClientCount       int `json:"client_count"`
	MessagesPerClient int `json:"messages_per_client"`
	ClientId          int `json:"client-id"`
}

func ConcurrentWorkerExecutor(totalWorkers int, maxExecutionsPerSecond int, fn workerFunc) {
	var wg sync.WaitGroup
	wg.Add(totalWorkers)

	start := time.Now()
	sem := make(chan int, maxExecutionsPerSecond)

	for i := 0; i < totalWorkers; i++ {
		go func(thingId int) {
			sem <- 1

			go func(thingId int) {
				defer wg.Done()

				if err := fn(thingId); err != nil {
					fmt.Println(err.Error())
					return
				}

				<-sem
			}(i)
		}(i)

		time.Sleep(time.Duration(1000/maxExecutionsPerSecond) * time.Millisecond)
	}

	wg.Wait()

	elapsed := time.Since(start)
	fmt.Printf("total execution time: %s\n", elapsed)
}
