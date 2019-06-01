package database

import "encoding/json"
import "bytes"
import "testing"
import "os"
import "parallels/database/common"


type BlobType int

const (
	Unclassified = iota
	Blob
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

	// insert Blob data to Store
	BlobID := "166507"
	rs := repoSet.Store.Upsert([]common.KeyValue{*common.NewKeyValue(Blob, BlobID, BlobData)})
	if !rs.IsSuccessful() {
		t.Error(rs.Error)
	}
	BlobID2 := "166508"
	rs = repoSet.Store.Upsert([]common.KeyValue{*common.NewKeyValue(Blob, BlobID2, BlobData)})
	if !rs.IsSuccessful() {
		t.Error(rs.Error)
	}

	// insert checkpoint 1 containing BlobIDs inserted to the Entity Store.
	ba, e := json.Marshal([]string{BlobID, BlobID2})
	if e != nil {
		t.Error(e)
	}

	rs = repoSet.NavigableStore.Upsert([]common.KeyValue{*common.NewKeyValue(Checkpoint, "1", ba)})
	if !rs.IsSuccessful() {
		t.Error(rs.Error)
	}

	r, rs := repoSet.NavigableStore.Get(Checkpoint, []string{"1"})
	if !rs.IsSuccessful() {
		t.Error(rs.Error)
	}

	var BlobIDs []string
	e = json.Unmarshal(r[0].Value, &BlobIDs)
	if e != nil {
		t.Error(e)
	}

	r, rs = repoSet.Store.Get(Blob, BlobIDs)
	if !rs.IsSuccessful() {
		t.Error(rs.Error)
	}
	if r == nil || len(r) != 2 {
		t.Errorf("Expected 2 Blobs not read.")
	}
	if r[0].Key != BlobID && r[0].Key != BlobID2 {
		t.Errorf("Expected 2 Blobs' keys not read.")
	}
	if r[1].Key != BlobID && r[1].Key != BlobID2 {
		t.Errorf("Expected 2 Blobs' keys not read.")
	}
	if !bytes.Equal(r[0].Value, r[1].Value) {
		t.Errorf("Expected Blobs retrieved from DB did not match.")
	}
}

var BlobData = []byte(`sample Blob with any data....`)
