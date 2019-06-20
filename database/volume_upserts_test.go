package database

import "fmt"
import "bytes"
import "testing"
import "os"
import "github.com/SharedCode/parallels/database/repository"
import "github.com/SharedCode/parallels/database/cassandra"

const count = 50000

func TestVolumeUpserts(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	var config, _ = LoadConfiguration(dir + "/config.json")
   repoSet, e := NewRepositorySet(config)
   defer cassandra.CloseSession()
	if e != nil {
		t.Error(e)
	}

   // NOTE: last ran using our ingestion test cassandra cluster took 8 mins to upsert 100,000 Albums.
   // This is about 3.5 hrs to upsert entire 30mil CDDB Albums.
   fmt.Printf("Performing %d upserts test...", count)

   const batchSize = 200
   batch := make([]repository.KeyValue, batchSize)
   ch := make(chan []repository.KeyValue)
   go sinker(ch, repoSet)

   for i := 1; i < count; i++ {
      //data := albumXML
      data := []byte("foobar")
      batch = append(batch, *repository.NewKeyValue(Album, fmt.Sprintf("%d",i), data))
      if len(batch) >= batchSize{
         ch <- batch
         batch = make([]repository.KeyValue, batchSize)
      }
   }
   // upsert the last part of the batch
   if len(batch) > 0{
      ch <- batch
   }
   close(ch)
   fmt.Println("Completed volume upserts, exitting.")
}

var index int
func sinker(ch chan []repository.KeyValue, repo repository.RepositorySet){
   for batch := range ch {
      index++
      go func(i int){
         rs := repo.Store.Set(batch...)
         if !rs.IsSuccessful() {
            fmt.Printf("Error Details:")
            fmt.Println(rs.ErrorDetails)
            fmt.Printf("Exitting test(batch#:%d), as failure occurred.", i)
            return
         }
      }(index)
   }
}

func TestVolumeReads(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	var config, _ = LoadConfiguration(dir + "/config.json")
	repoSet, e := NewRepositorySet(config)
	if e != nil {
		t.Error(e)
	}

   t.Logf("Performing %d reader test...", count)

   for i := 1; i < count; i++ {
      key := fmt.Sprintf("%d",i)
      albumData,rs := repoSet.Store.Get(Album, key)
      if !rs.IsSuccessful() {
         t.Error(rs.Error)
         t.Logf("Error Details:")
         t.Error(rs.ErrorDetails)
         t.Logf("Exitting test(loop:%d), as failure occurred.", i)
         return
      }
      if key != albumData[0].Key && !bytes.Equal(albumData[0].Value, albumXML){
         t.Errorf("Read data is not the same as expected (key:%s).", key)
         return
      }
   }
   t.Log("Completed volume reads, exitting.")
}
