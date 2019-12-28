package local

import "github.com/SharedCode/parallels"
import "github.com/SharedCode/parallels/database/repository"

// Pipeline interface defines the sourcer and sinker method signatures needing implementation.
// Pipeline interface defines the sourcer and sinker method signatures needing implementation.
type Pipeline interface {
	Sourcer() ([]repository.KeyValue,bool)
	Sinker([]repository.KeyValue)
}

// RunnerDefault is synonymous to Runner but using default parameter values
// on multiThreadedSourcer & threadCountThreshold.
func RunnerDefault(methods Pipeline){
	Runner(methods, false, parallels.DefaultThreadCountThreshold)
}

// Runner executes the pipeline runner.
func Runner(methods Pipeline, multiThreadedSourcer bool, threadCountThreshold int){
	// passthrough caller "cast" to expected types the params/returns.
	sourcer := func()(interface{},bool){
		return methods.Sourcer()
	}
	sinker := func(obj interface{}){
		b := obj.([]repository.KeyValue)
		methods.Sinker(b)
	}
	pi := parallels.NewParallel(threadCountThreshold)
	pi.Pipeline(sourcer, sinker, multiThreadedSourcer)
}
