package parallels

import "sync"

// DefaultThreadCountThreshold specifies a value that when count of threads running
// reaches this amount, ParallelFor will start to reduce aggressiveness in launching threads
// via go routine call. When count of threads become fewer than this amount, ParallelFor
// resumes aggressive thread launching. Your app can set this to a desired value.
// Default of 100 seems a good value, per how go runtime behaves on a commodity, 6 core CPU PC.
var DefaultThreadCountThreshold = 100

// ParallelInfo interface defines the sourcer and sinker method signatures needing implementation.
type ParallelInfo struct {
	ThreadCountThreshold int
}

// NewParallelDefault is synonymous to NewParallel but setting ThreadCountThreshold to the default value.
func NewParallelDefault() *ParallelInfo {
	return NewParallel(DefaultThreadCountThreshold)
}

// NewParallel instantiates a new Parallel object. threadCountThreshold is expected to be > 0,
// otherwise it is set to 1.
func NewParallel(threadCountThreshold int) *ParallelInfo {
	// apply some reasonable bounds on thread counts threshold.
	if threadCountThreshold < 1 {
		threadCountThreshold = 1
	}
	const MaxThreadCountThreshold = 100000
	if threadCountThreshold > DefaultThreadCountThreshold && threadCountThreshold > MaxThreadCountThreshold {
		threadCountThreshold = DefaultThreadCountThreshold
	}

	return &ParallelInfo{
		ThreadCountThreshold: threadCountThreshold,
	}
}

// Pipeline invokes sourcer in a for-loop in a dedicated thread(or multithreaded if 'multiThreadSourcing' is true),
// sourcer should return 'nil,true' when done. Items sourced from the call are then 'sinked' to the channel.
// In which, the sinker function will receive, process these items in parallel.
func (pi ParallelInfo) Pipeline(sourcer func() (interface{}, bool), sinker func(interface{}), multiThreadSourcing bool) {
	if !multiThreadSourcing {
		ch := make(chan interface{})
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			for {
				obj, done := sourcer()
				if done {
					if obj != nil {
						ch <- obj
					}
					break
				}
				ch <- obj
			}
			close(ch)
			wg.Done()
		}()
		pi.ParallelFor(ch, sinker, &wg)
		wg.Wait()
		return
	}

	// multi-threaded sourcing...
	var ctr int
	threadCountThreshold := pi.ThreadCountThreshold / 2
	if threadCountThreshold < 2 {
		threadCountThreshold = 2
	}
	ch := make(chan interface{})
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		var sourcingDone bool
		for !sourcingDone {
			sourcerFunc := func() {
				// short circuit if we're done (other threads may signal 'done')
				if sourcingDone {
					return
				}
				obj, done := sourcer()
				if done {
					sourcingDone = true
					if obj != nil {
						ch <- obj
					}
					return
				}
				ch <- obj
			}
			f := func() {
				defer func() {
					wg.Done()
					ctr--
				}()
				sourcerFunc()
			}
			if ctr >= threadCountThreshold {
				//fmt.Printf("source sync call after %d threads...\n", ctr)
				sourcerFunc()
				continue
			}
			wg.Add(1)
			ctr++
			go f()
		}
		close(ch)
		wg.Done()
	}()
	parallelFor(ch, sinker, &wg, threadCountThreshold)
	wg.Wait()
}

// ParallelFor does a parallelized for loop on a given channel and invokes
// 'action' function on each item received from the channel.
// On each iteration pass, 'wg.Done()' is invoked to signal completion of the 'action'
// worker thread, to the caller.
//
// This ParallelFor function does a simple logic to prevent thread over-allocation.
// Which is, one of its (key/) primary uses.
func (pi ParallelInfo) ParallelFor(sourceChannel <-chan interface{}, targetAction func(interface{}), wg *sync.WaitGroup) {
	parallelFor(sourceChannel, targetAction, wg, pi.ThreadCountThreshold)
}

func parallelFor(sourceChannel <-chan interface{}, targetAction func(interface{}), wg *sync.WaitGroup, threadCountThreshold int) {
	var ctr int
	if threadCountThreshold < 0 {
		threadCountThreshold = 0
	}
	wg.Add(1)
	for obj := range sourceChannel {
		f := func(obj interface{}) {
			defer func() {
				wg.Done()
				ctr--
			}()
			targetAction(obj)
		}

		//fmt.Printf("thread count %d.\n", ctr)

		if ctr >= threadCountThreshold {
			// perform synchronous call when threshold count of threads is reached
			// to prevent thread over-allocation.
			//fmt.Printf("sync call after %d threads...\n", ctr)
			targetAction(obj)
			continue
		}
		wg.Add(1)
		ctr++
		go f(obj)
	}
	wg.Done()
}
