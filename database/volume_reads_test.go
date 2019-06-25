package database

import "fmt"
import "testing"
import "sync"
import "time"
import "os"
import "github.com/SharedCode/parallels/database/repository"

func TestVolumeReads(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	var config, _ = LoadConfiguration(dir + "/config.json")
   repoSet, e := NewRepositorySet(config)
   defer CloseSession()
	if e != nil {
		t.Error(e)
	}

   fmt.Printf("Performing %d reads test...\n", count)

   const batchSize = 1000
   batch := make([]string, 0, batchSize)
   ch := make(chan []string)
   var wg sync.WaitGroup
   go func(){
      for i := 1; i < count; i++ {
         batch = append(batch, fmt.Sprintf("%d",i))
         if len(batch) >= batchSize{
            wg.Add(1)
            ch <- batch
            batch = make([]string, 0, batchSize)
         }
      }
      // upsert the last part of the batch
      if len(batch) > 0{
         wg.Add(1)
         ch <- batch
      }
      close(ch)
   }()

   reader(ch, repoSet, &wg)
   wg.Wait()
   fmt.Println("Completed volume upserts, exitting.")
}

func reader(ch chan []string, repo repository.RepositorySet, wg *sync.WaitGroup) {
   for batch := range ch {
      index++
      f := func(i int, keys []string){
         defer wg.Done()
         for i2 := 0; i2 < 10; i2++ {
            kv,rs := repo.Store.Get(Album, keys...)
            if rs.IsSuccessful() {
               if len(kv) != len(keys){
                  fmt.Printf("Failed, expected item count %d, got %d\n", len(keys), len(kv))
               } else {
                  fmt.Printf("Successful read batch# %d\n", i)
               }
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
      if index%25 == 0{
         // perform synchronous call every 25(!) threads to prevent thread over-allocation.
         // hacky way of thread mgmt in a "test" script. :)
         f(index, batch)
         continue
      }
      go f(index, batch)
   }
}
