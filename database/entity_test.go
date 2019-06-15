package database

import "testing"
import "bytes"
import "os"
import "encoding/json"
import "github.com/SharedCode/parallels/database/repository"

type AlbumType int

const (
	Unclassified = iota
	Album
	Checkpoint
)

func TestBasic(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	var config, _ = LoadConfiguration(dir + "/config.json")
	repoSet, e := NewRepositorySet(config)
	if e != nil {
		t.Error(e)
	}

	// insert Album XML w/ ID 1665078 to Store
	albumID := "1665078"
	rs := repoSet.Store.Set(*repository.NewKeyValue(Album, albumID, albumXML))
	if !rs.IsSuccessful() {
		t.Error(rs.Error)
	}
	albumID2 := "1665079"
	rs = repoSet.Store.Set(*repository.NewKeyValue(Album, albumID2, albumXML))
	if !rs.IsSuccessful() {
		t.Error(rs.Error)
   }
   
   ba, e := json.Marshal([]string{albumID, albumID2})
 	if e != nil {
 		t.Error(e)
    }
    
	// insert checkpoint 1 containing AlbumIDs inserted to the Entity Store.
	rs = repoSet.NavigableStore.Set(*repository.NewKeyValue(Checkpoint, "1", ba))
	if !rs.IsSuccessful() {
		t.Error(rs.Error)
	}

	gr,rs := repoSet.NavigableStore.Get(Checkpoint, "1")
	if !rs.IsSuccessful() {
		t.Error(e)
   }

   var albumIDs []string
	e = json.Unmarshal(gr[0].Value, &albumIDs)
	if e != nil {
		t.Error(e)
	}
	gr,rs = repoSet.Store.Get(Album, albumIDs...)
	r := gr
	if !rs.IsSuccessful() {
		t.Error(rs.Error)
	}
	if r == nil || len(r) != 2 {
		t.Errorf("Expected 2 Albums not read.")
	}
	if r[0].Key != albumID && r[0].Key != albumID2 {
		t.Errorf("Expected 2 Albums' keys not read.")
	}
	if r[1].Key != albumID && r[1].Key != albumID2 {
		t.Errorf("Expected 2 Albums' keys not read.")
	}
   if !bytes.Equal(r[0].Value, r[1].Value) {
		t.Errorf("Expected AlbumXML retrieved from DB did not match.")
	}
}

var albumXML = []byte(`Blob data`)
