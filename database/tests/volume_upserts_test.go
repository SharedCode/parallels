package tests

import "fmt"
import "sync"
import "time"
import "testing"
import "os"
import "github.com/SharedCode/parallels/database/repository"
import "github.com/SharedCode/parallels/database"
import "github.com/SharedCode/parallels"

const count = 100000

// High volume upserts "test"! 'the DB passed with flying colors, per "performance" expectations. :)

func TestVolumeUpserts(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	var config, _ = database.LoadConfiguration(dir + "/config.json")
   repoSet, e := database.NewRepositorySet(config)
   defer database.CloseSession()
	if e != nil {
      t.Error(e)
      t.FailNow()
	}

   fmt.Printf("Performing %d upserts test...\n", count)

   const batchSize = 1000
   batch := make([]repository.KeyValue, 0, batchSize)
   ch := make(chan []repository.KeyValue)
   var wg sync.WaitGroup
   go func(){
      for i := 1; i < count; i++ {
         data := albumXML
         batch = append(batch, *repository.NewKeyValue(Album, fmt.Sprintf("%d",i), data))
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

   writer(ch, repoSet, &wg)
   wg.Wait()
   fmt.Println("Completed volume upserts, exitting.")
}

func writer(ch chan []repository.KeyValue, repo repository.RepositorySet, wg *sync.WaitGroup) {
   var ctr int
   var index int
   for batch := range ch {
      ctr++
      index++
      f := func(i int, b []repository.KeyValue){
         defer func(){
            wg.Done()
            ctr--
         }()
         for i2 := 0; i2 < 10; i2++ {
            rs := repo.Store.Set(b...)
            if rs.IsSuccessful() {
               fmt.Printf("Successful upsert batch# %d\n", i)
               return
            }
            if rs.ErrorDetails != nil{
               fmt.Printf("Error #%d, msg: %s", i, rs.Error)
               fmt.Printf("Error Details:")
               fmt.Println(rs.ErrorDetails)
            }
            time.Sleep(4*time.Second)
         }
         fmt.Printf("Error persisted for 10 times, giving up\n")
      }
      if ctr > parallels.DefaultThreadCountThreshold{
         // perform synchronous call every 25(!) threads to prevent thread over-allocation.
         // hacky way of thread mgmt in a "test" script. :)
         f(index, batch)
         continue
      }
      go f(index, batch)
   }
}
