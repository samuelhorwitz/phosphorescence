package spotifyclient

import (
	"sync"
)

// Go WaitGroups have a restriction that you cannot use `Wait`
// until after `Add` is called when the counter is zero. This
// limits it's use to times you are definitely adding deltas
// and then subsequently waiting for completetion, even though
// what we need is something that opens and closes as needed
// while everywhere else blocks until a broadcast wakes it up.
type waitGroupCond struct {
	cond  *sync.Cond
	count int64
}

func (wg *waitGroupCond) maybeInit() {
	if wg.cond == nil {
		wg.cond = &sync.Cond{L: &sync.Mutex{}}
	}
}

func (wg *waitGroupCond) Increment() {
	wg.Add(1)
}

func (wg *waitGroupCond) Add(delta int64) {
	wg.maybeInit()
	wg.cond.L.Lock()
	wg.count += delta
	if wg.count < 0 {
		panic("negative waitGroupCond counter")
	} else if wg.count == 0 {
		wg.cond.Broadcast()
	}
	wg.cond.L.Unlock()
}

func (wg *waitGroupCond) Done() {
	wg.Add(-1)
}

func (wg *waitGroupCond) Wait() {
	wg.maybeInit()
	wg.cond.L.Lock()
	for wg.count != 0 {
		wg.cond.Wait()
	}
	wg.cond.L.Unlock()
}
