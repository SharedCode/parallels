package parallels


import "testing"
import "os"

import "fmt"
import "encoding/json"
import "github.com/SharedCode/parallels/database"

const (
	sag = "TestScatterAndGather"
)

func TestScatterAndGather(t *testing.T){
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	var config, _ = database.LoadConfiguration(dir + "/config.json")
   repoSet, e := database.NewRepositorySet(config)
   defer database.CloseSession()
	if e != nil {
		t.Error(e)
	}

	runner := NewRunner(repoSet.Store, processor);

	var jobs = []int{1,2,3,4,5}
	ba,_ := json.Marshal(jobs)
	_,e = runner.ScatterAndGather([][]byte{ba}, sag)
	if e == nil{
		return
	}
	t.Error(e)
}

func processor(items []byte, actionName string) error{
	if actionName == sag{
		var jobs []int
		e := json.Unmarshal(items, &jobs)
		if e != nil{
			return e
		}
		return processJob(jobs)
	}
	return nil
}

func processJob(jobs interface{}) error{
	localFunc := func(j int){
		fmt.Printf("Job %d was processed!\n", j)
	}
	if t,ok := jobs.(int); ok{
		localFunc(t)
	}
	if t,ok := jobs.([]int); ok{
		for _,j := range t{
			localFunc(j)
		}
		return nil
	}
	return fmt.Errorf("Can't deserialize received 'jobs' data")
}
