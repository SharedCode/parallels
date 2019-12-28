package parallels

//import "github.com/SharedCode/parallels/database/cache"
import "github.com/SharedCode/parallels/database/repository"

type Runner struct{
	Store repository.Repository
	Dispatcher Processor
}

type ActionResult struct{
	ActionName string
	ActionID string
	Result []byte
	Error error
}

func (r ActionResult) IsSuccessful() bool{
	return r.Error == nil
}

// JobDistributor contains logic to calculate how many Jobs are to be assigned per Node in the Cluster.
// Default implementation simply, will "enqueue" (to Redis) the Jobs with no Node assignment information of any kind.
type JobDistributor interface{
}
// JobSourcer contains logic to "take" Jobs from the queue (Redis) and mark these Jobs as appropriate.
// See JobStatus for different Job State that the JobSource will set.
type JobSourcer interface{
}

type Processor func(actionData []byte, actionName string) error

// NewRunner creates a new Runner instance.
func NewRunner(c repository.Repository, dispatcher Processor) Runner{
	// assign an in-memory cache if cache is nil. todo: finalize whether we need this!
	if c == nil {
		//c = cache.NewMemoryCache()
	}
	return Runner{
		Store : c,
		Dispatcher: dispatcher,
	}
}


func (runner Runner) ScatterAndGather(batchActionData [][]byte, actionName string) (chan ActionResult,error){
	results,e := runner.Scatter(batchActionData, actionName)
	if e != nil{
		return nil,e
	}

	var bids []string
	for i := range results{
		bids = append(bids, results[i].ActionID)
	}
	return runner.Gather(bids, actionName)
}

func (runner Runner) WaitUntilGathered(resultCount int, ch chan ActionResult) ([]ActionResult, error){
	r2 := make([]ActionResult, 0, resultCount)
	for i := 0; i < resultCount; i++{
		result := <- ch
		r2 = append(r2, result)
	}
	return r2,nil
}

// Scatter will distribute processing to (potentially) many workers locally and remotely (in the cluster).
func (runner Runner) Scatter(batchActionData [][]byte, actionName string) ([]ActionResult, error){
	return nil,nil
}

// Gather will collect result(s) of the processing done to each of the action as denoted by its ID (actionID).
// Each of the ActionResult, as it arrives or detected, will be "sinked" to the channel returned.
func (runner Runner) Gather(batchActionID []string, actionName string) (chan ActionResult,error){
	
	return nil,nil
}

// RunSomeActions will get jobs from repo(e.g. - Redis) then dispatch them to the processor function(Dispatcher) 
// received as parmeter during construction of the Runner object, see NewRunner dispatcher parameter.
// NOTE: this method is designed for call on the background during idle time, so the running App can do its share
// of processing workloads in behalf of the "app cluster".
func (runner Runner) RunSomeActions(actionCount int) ([]ActionResult, error){

	// Use JobSourcer to get Jobs and to manage these Jobs as they get processed...
	// Call Dispatcher for each Job gotten from Queue.
	// Mark Job(s) that got (successfully) completed so they don't get re-processed again.

	return nil,nil
}
