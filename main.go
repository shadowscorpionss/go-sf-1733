package main

import (
	"fmt"
	"sync"
)

// Конечное значение счетчика
const endCounterValue int64 = 1000

// количество потоков
const routineCount int64 = 37

// Шаг наращивания счётчика
const step int64 = endCounterValue / routineCount

func main() {

	var counter int64 = 0
	var wg sync.WaitGroup

	c := sync.NewCond(&sync.Mutex{})

	fmt.Println("starting", routineCount+1, "routines to count and check", endCounterValue)

	increment := func(step int64) {
		defer wg.Done()
		c.L.Lock()
		c.Wait()
		counter += step
		c.L.Unlock()
	}

	check := func() {
		defer wg.Done()
		for {
			c.L.Lock()
			if counter >= endCounterValue {
				break
			}
			c.L.Unlock()

			c.Broadcast()
		}
	}

	wg.Add(1)
	go check()

	for i := int64(0); i < routineCount; i++ {
		wg.Add(1)
		go increment(step)
	}
	if routineCount*step < endCounterValue {
		fmt.Println("adding one routine to complete")

		wg.Add(1)
		go increment(endCounterValue - routineCount*step)
	}

	// Ожидаем поступления сигнала
	wg.Wait()
	// Печатаем результат, надеясь, что будет 1000
	fmt.Println("finished. our result is", counter)
}
