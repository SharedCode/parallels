package parallels

import "fmt"
import "testing"
import "sync"
import "time"
import "os"
import "github.com/SharedCode/parallels/database"
import "github.com/SharedCode/parallels/database/repository"

const count = 9999999

// High volume upserts "test"! 'the DB passed with flying colors, per "performance" expectations. :)

func TestParallelFor(t *testing.T) {
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

   fmt.Printf("Performing %d upserts test...\n", count)

   const batchSize = 1000
   batch := make([]repository.KeyValue, 0, batchSize)
   ch := make(chan interface{})
   var wg sync.WaitGroup
   go func(){
      for i := 1; i < count; i++ {
         data := "foobarGroup"
         batch = append(batch, *repository.NewKeyValue(data, fmt.Sprintf("%d",i), []byte(data)))
         if len(batch) >= batchSize{
            wg.Add(1)
            ch <- batch
            batch = make([]repository.KeyValue, 0, batchSize)
         }
      }
      // upsert the last part of the batch
      if len(batch) > 0{
         wg.Add(1)
         ch <- batch
      }
      close(ch)
   }()

   f := func(obj interface{}){
      b := obj.([]repository.KeyValue)
      for i2 := 0; i2 < 10; i2++ {
         rs := repoSet.Store.Set(b...)
         if rs.IsSuccessful() {
            fmt.Printf("Successful upsert batch w/ 1st item key %s\n", b[0].Key)
            return
         }
         if rs.ErrorDetails != nil{
            fmt.Printf("Error msg: %s", rs.Error)
            fmt.Printf("Error Details:")
            fmt.Println(rs.ErrorDetails)
         }
         time.Sleep(4*time.Second)
      }
      fmt.Printf("Error persisted for 10 times, giving up\n")
   }

   ParallelFor(ch, f, &wg)
   wg.Wait()
   fmt.Println("Completed volume upserts, exitting.")
}

func TestParallelFor(t *testing.T) {
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

   fmt.Printf("Performing %d upserts test...\n", count)

   const batchSize = 1000

   sourcer := func()(func()(interface{},bool)){
      ctr := 0
      return func()(interface{},bool){
         ctr++
         if ctr > count{
            return nil,true
         }
         batch := make([]repository.KeyValue, 0, batchSize)
         for i := 0; i < batchSize; i++ {
            data := "foobarGroup"
            batch = append(batch, *repository.NewKeyValue(data, fmt.Sprintf("%d",i), []byte(data)))
         }
         return batch,false
      }
   }

   sinker := func(obj interface{}){
      b := obj.([]repository.KeyValue)
      for i2 := 0; i2 < 10; i2++ {
         rs := repoSet.Store.Set(b...)
         if rs.IsSuccessful() {
            fmt.Printf("Successful upsert batch w/ 1st item key %s\n", b[0].Key)
            return
         }
         if rs.ErrorDetails != nil{
            fmt.Printf("Error msg: %s", rs.Error)
            fmt.Printf("Error Details:")
            fmt.Println(rs.ErrorDetails)
         }
         time.Sleep(4*time.Second)
      }
      fmt.Printf("Error persisted for 10 times, giving up\n")
   }
   SourceSink(sourcer(), sinker)
   fmt.Println("Completed volume upserts, exitting.")
}
