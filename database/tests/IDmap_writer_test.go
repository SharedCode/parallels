package tests

import "fmt"
import "strings"
import "time"
import "sync"
import "testing"
import "strconv"
import "os"
import "github.com/SharedCode/parallels/database/repository"
import "github.com/SharedCode/parallels/database"
import "github.com/SharedCode/parallels"

const itemcount = 10000
const maxPartitionCount = 100

func TestIDMapWriter(t *testing.T) {
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

   fmt.Printf("Performing %d ID map writer test...\n", itemcount)
   const batchSize = 1000

   sourcer := func()(func()(interface{},bool)){
      ctr := 0
      var ended bool
      return func()(interface{},bool){
         if ended {
            return nil,true
         }
         ctr++
         fmt.Printf("sourcer batch # %d.\n", ctr)
         batch := make([]repository.KeyValue, 0, batchSize)
         for i := 1; i < batchSize && !ended; i++ {
            if (ctr-1)*batchSize+i > itemcount{
               ended = true
               if batch == nil || len(batch) == 0 {return nil,true}
            }
            data := albumXML
            batch = append(batch, *repository.NewKeyValue(Album, fmt.Sprintf("%d",(ctr-1)*batchSize+i), data))
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
            // write the data group and the album IDs for the batch.
            idMapWriter(b, repoSet)
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
      // intentionally not handling error here, as not the focus of the test.
   }
   parallel := parallels.NewParallelDefault()
   parallel.Pipeline(sourcer(), sinker, false)

   fmt.Println("Completed time series writer test, exitting.")
}

func idMapWriter(b []repository.KeyValue, repo repository.RepositorySet){
   sb := strings.Builder{}
   for _,v := range b {
      if sb.Len() > 0 {
         sb.WriteString(",")
      }
      sb.WriteString(v.Group + "_" + v.Key)
   }
   // use modulo to spread out evenly, through different partitions.
   partitionKey := getCurrentTime().Unix() % maxPartitionCount
   k := fmt.Sprintf("%s_%s", b[0].Group, b[0].Key)
   fmt.Printf("Writing to partition #%d, key %s.\n", partitionKey, k)

   // assumes batch is for the same data group!
   // populate DB targeting calculated partition ID & generated entity IDs of the batch.
   repo.NavigableStore.Set(*repository.NewKeyValue(strconv.FormatInt(partitionKey, 10), k, []byte(sb.String())))
}


func TestIDPartitionGen(t *testing.T) {
   for i := 0; i < 200; i++ {
      partitionKey := getCurrentTime().Unix() % maxPartitionCount
      fmt.Printf("Writing to partition #%d.\n", partitionKey)
   }
}

func getCurrentTime() time.Time {
   ts_locker.Lock()
   defer ts_locker.Unlock()

   // add 1 minute to simulate call every minute.
   simulated_current_time = simulated_current_time.Add(1*time.Second)
   return simulated_current_time
}

var ts_locker sync.Mutex
var simulated_current_time = time.Now().UTC()
