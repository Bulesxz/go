package time

import (
	"sync"
	"time"
	"container/list"
	"fmt"
)

type TimingWheel struct {
	sync.Mutex

	interval time.Duration  //间隔

	ticker *time.Ticker
	quit   chan struct{}

	buckets int //桶的个数

	tasks    []*list.List
	
	pos int
}

func NewTimingWheel(interval time.Duration, buckets int) *TimingWheel {
	w := new(TimingWheel)

	w.interval = interval

	w.quit = make(chan struct{})
	w.pos = 0

	w.buckets = buckets
		
	w.tasks = make([]*list.List,buckets)
	for i:=0;i<buckets;i++{
		w.tasks[i]=list.New()
		//fmt.Println("---",i,w.tasks[i])
	}
	
	w.ticker = time.NewTicker(interval)
	
	return w
}

func (w *TimingWheel) Stop() {
	close(w.quit)
}

func (w *TimingWheel) Add(timeout time.Duration,callback func()) {
	if timeout >=  w.interval * time.Duration(w.buckets) {
		panic("timeout too much, over maxtimeout")
	}

	w.Lock()

	index := (w.pos + int(timeout/w.interval)) % w.buckets
	w.tasks[index].PushBack(callback)
	//fmt.Println("++++",index,w.tasks[index])
	
	w.Unlock()

}

func (w *TimingWheel) Run() {
	for {
		select {
		case <-w.ticker.C:
			w.onTicker()
		case <-w.quit:
			w.ticker.Stop()
			return
		}
	}
}

func (w *TimingWheel) onTicker() {
	w.Lock()
	
	
	//fmt.Println("pos=",w.pos,w.tasks[w.pos],w.tasks[w.pos].Front())
	
	task := w.tasks[w.pos]
	w.tasks[w.pos]=nil
	w.pos = (w.pos + 1) % w.buckets
	w.Unlock()
	
	if task != nil {
		//fmt.Println("task len",w.tasks[w.pos])
		doTask(task)
	}
}

func doTask(l *list.List){
	
	for e := l.Front(); e != nil; e = e.Next() {
		callback :=e.Value.(func())
        //fmt.Println("datask",e.Value.(func()))
		//task.Front().Value.(func())
		callback()
    }
	
}


func timeOut(){
	fmt.Println("timeOut")
}
/*
func main() {
	w := NewTimingWheel(1*time.Second, 10)
	
	w.Add(5 * time.Second,timeOut)
	
	time.Sleep(1000*time.Second) 
}*/