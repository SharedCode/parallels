package parallels

import "sync"

// Run will create an invoke wrapper that will call Sourcer
// and Sinker functions in parallel.
// This version doesn't do "remote calls".
func Run(sourcer func() (interface{},bool), sinker func(interface{})) {
	ch := make(chan interface{})
	var wg sync.WaitGroup

	go func(){
		for{
			wg.Add(1)
			obj,done := sourcer()
			ch <- obj
			if done {
				break
			}
		}
		close(ch)
	}()
	sink(ch, sinker, &wg)
	wg.Wait()
}

func sink(ch chan interface{}, sinker func(obj interface{}), wg *sync.WaitGroup){
	var ctr int
	for obj := range ch {
		ctr++
		f := func(obj interface{}){
		   defer wg.Done()
		   sinker(obj)
		}

		if ctr%25 == 0{
		   // perform synchronous call every iteration set(!). This prevents thread over-allocation.
		   f(obj)
		   continue
		}
		go f(obj)
	 }
}
