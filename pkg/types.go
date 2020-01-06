package awsiotloadsimulator

import (
	"fmt"
	"sync"
	"time"
)

type workerFunc func(int) error

type SimulationRequest struct {
	StartClientNumber      int `json:"start_client_num"`
	ClientCount            int `json:"client_count"`
	MessagesPerClient      int `json:"messages_per_client"`
	SecondsBetweenMessages int `json:"seconds_between_messages"`
	ClientId               int `json:"client-id"`
}

func (s SimulationRequest) String() string {
	return fmt.Sprintf(`SimulationRequest
		StartClientNumber: %d
		ClientCount: %d
		MessagesPerClient: %d
		SecondsBetweenMessages: %d
		ClientId: %d`,
		s.StartClientNumber,
		s.ClientCount,
		s.MessagesPerClient,
		s.SecondsBetweenMessages,
		s.ClientId,
	)
}

func ConcurrentWorkerExecutor(totalWorkers int, maxExecutionsPerSecond time.Duration, fn workerFunc) time.Duration {
	var wg sync.WaitGroup
	wg.Add(totalWorkers)

	start := time.Now()
	// sem := make(chan int, maxExecutionsPerSecond)

	rate := time.Second / maxExecutionsPerSecond
	throttle := time.Tick(rate)

	for i := 0; i < totalWorkers; i++ {
		<-throttle // rate limit our Service.Method RPCs

		// go func(thingId int) {
		// 	sem <- 1

		go func(thingId int) {
			defer wg.Done()

			if err := fn(thingId); err != nil {
				fmt.Println(err.Error())
				return
			}

			// <-sem
		}(i)
		// }(i)
	}

	wg.Wait()

	elapsed := time.Since(start)
	return elapsed
}
