package base

import (
	"sync"
	"sync/atomic"
	"time"
)

func BenchmarkFunc(n int32, taskthread int, f func() bool) (float64, int32) {
	var failNum int32 = 0
	var usetime float64 = 0
	var wg sync.WaitGroup

	start := time.Now().UnixNano()
	for i := 0; i < taskthread; i++ {
		wg.Add(1)
		go func(i int) {
			//fmt.Println("go add ", i)
			for atomic.AddInt32(&n, -1) >= 0 {
				fail := f()
				if fail == false {
					atomic.AddInt32(&failNum, 1)
				}
			}
			wg.Done()
			//fmt.Println("go done ", i)
		}(i)
	}

	wg.Wait()
	end := time.Now().UnixNano()
	usetime = float64((end - start) / 1000000) //mill

	return usetime, failNum
}
