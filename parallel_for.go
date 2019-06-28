package parallels

import "sync"
import "fmt"

// ThreadCountThreshold specifies a value that when count of threads running
// reaches this amount, ParallelFor will start to reduce aggressiveness in launching threads
// via go routine call. When count of threads become fewer than this amount, ParallelFor
// resumes aggressive thread launching. Your app can set this to a desired value.
// Default of 100 seems a good value, per how go runtime behaves on a commodity, 6 core CPU PC.
var ThreadCountThreshold = 100

// PipelineParallelForSink invokes sourcer in a for loop in a dedicated single thread.
// Items sourced from the call are then 'sinked' to the channel.
// In which, the sinker function will receive, process these items in parallel for loop.
func PipelineParallelForSink(sourcer func() (interface{},bool), sinker func(interface{})){
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

// ParallelForPipeline does a sourcing-sinking run. It will invoke 'sourcer' function
// in a parallel for loop then each item it returns, will be written to the channel
// that which, will be passed on to the 'sinker' function, which also runs in parallel for
// listening on the channel.
//
// Running these 'sourcing' and 'sinking' functions in parallel threads and
// communicating via a channel as medium, are all done in this function,
// so your code/function can focus on tasks at hand, i.e. - to source and to sink
// domain logic.
func ParallelForPipeline(sourcer func() (interface{},bool), sinker func(interface{})){
	var ctr int
	ch := make(chan interface{})
	var wg sync.WaitGroup
	go func(){
		var sourcingDone bool
		for !sourcingDone {
			wg.Add(1)
			ctr++
			f := func(){
				obj,done := sourcer()
				ctr--
				if done {
					sourcingDone = true
					if obj != nil{
						ch <- obj
					}
					return
				}
				ch <- obj
			}
			if ctr > ThreadCountThreshold {
			   fmt.Printf("source sync call after %d threads...\n", ctr)
			   f()
			   continue
			}
			go f()
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
func ParallelFor(sourceChannel chan interface{}, targetAction func(interface{}), wg *sync.WaitGroup) {
	var ctr int
	for obj := range sourceChannel {
		ctr++
		f := func(obj interface{}){
		   defer func(){
			wg.Done()
			ctr--
		   }()
		   targetAction(obj)
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
