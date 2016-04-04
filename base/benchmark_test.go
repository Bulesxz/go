package base

import (
	"fmt"
	"testing"
	"time"
)

func Test_benchmark(t *testing.T) {
	f := func() bool {
		fmt.Println("test")
		time.Sleep(1 * time.Second)
		return false
	}
	usetime, failNum := BenchmarkFunc(10, 1, f)
	fmt.Println(usetime, failNum)
}
