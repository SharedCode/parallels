package parallels

import "sync"
import "fmt"

var ThreadCountThreshold = 500

// SourceSink does a pipeline run. It will run 'sourcer'
// in a for loop then each item it returns, will be written to the channel
// that which, will be passed on to the 'sinker' function.
// Running these 'sourcing' and 'sinking' functions in parallel threads and
// communicating via a channel as medium, are all done in this function,
// so your code/function can focus on tasks at hand, i.e. - to source and to sink
// domain logic.
func SourceSink(sourcer func() (interface{},bool), sinker func(interface{})){
	ch := make(chan interface{})
	var wg sync.WaitGroup
	go func(){
		for {
			wg.Add(1)
			obj,done := sourcer()
			if done {
				if obj != nil{
					ch <- obj
				}
				break
			}
			ch <- obj
		}
		close(ch)
	}()
	ParallelFor(ch, sinker, &wg)
	wg.Wait()
}

// ParallelFor does a parallelized for loop on a given channel and invokes
// 'action' function on each item received from the channel.
// On each iteration pass, 'wg.Done()' is invoked to signal completion of the 'action'
// worker thread, to the caller.
//
// This ParallelFor function does a simple logic to prevent thread over-allocation.
// Which is, one of its (key/) primary uses.
func ParallelFor(ch chan interface{}, action func(interface{}), wg *sync.WaitGroup) {
	var ctr int
	for obj := range ch {
		ctr++
		f := func(obj interface{}){
		   defer func(){
			wg.Done()
			ctr--
		   }()
		   action(obj)
		}

		fmt.Printf("thread count %d.\n", ctr)

		if ctr > ThreadCountThreshold {
		   // perform synchronous call when threshold count of threads is reached
		   // to prevent thread over-allocation.
		   fmt.Printf("sync call after %d threads...\n", ctr)
		   f(obj)
		   continue
		}
		go f(obj)
	 }
}
