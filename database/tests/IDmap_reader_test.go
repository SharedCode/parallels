package tests

import "fmt"
import "sync"
import "testing"
import "strconv"
import "os"
import "github.com/SharedCode/parallels/database/repository"
import "github.com/SharedCode/parallels/database"
import "github.com/SharedCode/parallels"

func TestIDMapReader(t *testing.T) {
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

   fmt.Printf("Performing %d ID map reader test...\n", itemcount)
   const batchSize = 1000

   var wg sync.WaitGroup
   ch := make(chan interface{})
   wg.Add(1)
   go func(){
      defer func(){
         close(ch)
         wg.Done()
      }()
      for i := int64(0); i < maxPartitionCount; i++{
         fmt.Printf("sourcer: partition # %d.\n", i)

         filter := repository.Filter{MaxCountLimit: 100}
         for {
            v,r := repoSet.NavigableStore.Navigate(strconv.FormatInt(i,10), filter)
            if !r.IsSuccessful(){
               fmt.Printf("Failed reading partition # %d.", i)
               return;
            }
            if v == nil || len(v) == 0{
               fmt.Printf("Done reading contents of partition # %d.", i)
               break;
            }
            ch <- v
            filter.UpperboundKey = v[len(v)-1].Key
         }
      }
   }()

   albumReader := func(obj interface{}){
      b := obj.([]repository.KeyValue)
      for i := 0; i < len(b); i++ {
         fmt.Printf("About to read data of partition #%s, key %s.\n", b[i].Group, b[i].Key)

         // todo: parse b[i].Key into group and key, e.g. - Entity Group ("Album") and 
         // AlbumNumber for use to read Album record.

      }
   }
   parallel := parallels.NewParallelDefault()
   parallel.ParallelFor(ch, albumReader, &wg)
   wg.Wait()

   fmt.Println("Completed ID map reader test, exitting.")
}
