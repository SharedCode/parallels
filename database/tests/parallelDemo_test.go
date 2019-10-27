package tests

import "fmt"
import "testing"
import "github.com/SharedCode/parallels/database/repository"
import "github.com/SharedCode/parallels"

// TestParallelDemo is a simple test to demonstrate how to use Parallel Pipeline runner API.
func TestParallelDemo(t *testing.T) {
	fmt.Println("Performing test parallel demo.")
	d := newDemo()
	parallels.RunDefault(d)
}

type demoStruct struct {
	batchSize      int
	batchCount     int
	currentBatchNo int
}

func newDemo() parallels.Pipeline {
	return &demoStruct{batchSize: 1000, batchCount: 10}
}

// Sourcer method simulates reading from a source and returning the data, in batches!
// NOTE: it is recommended to operate in batches of data so an implementation can take advantage of "bulk"
// I/O or operations available in the Sourcing and Sinking sides. Most/all databases or storage
// devices, processes, etc... operate & can optimize if given a set of data, instead of an item each call.
func (d *demoStruct) Sourcer() ([]repository.KeyValue, bool) {
	if d.currentBatchNo >= d.batchCount {
		// return nil and true to imply that we are done sourcing.
		return nil, true
	}
	batch := make([]repository.KeyValue, d.batchSize)
	for i := 0; i < d.batchSize; i++ {
		batch[i] = *repository.NewKeyValue(fmt.Sprintf("g%d", d.currentBatchNo),
			fmt.Sprintf("k%d", d.currentBatchNo*d.batchSize+i), []byte(fmt.Sprintf("v%d", d.currentBatchNo*d.batchSize+i)))
	}
	d.currentBatchNo++
	return batch, false // return the batch and false to imply there is "more" data to be sourced.
}

// Sinker method simulates writing to a target a batch of received data (from sourcer).
func (d *demoStruct) Sinker(kvps []repository.KeyValue) {
	// since this is invoked in parallel, data will arrive here NOT in order,
	// which is as expected. In real world situation, this batch of data can be written out to a database, or for output, etc...
	for _, k := range kvps {
		fmt.Printf("Sinking group:%s, key:%s, value:%s\n", k.Group, k.Key, string(k.Value))
	}
}
