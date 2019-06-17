package database

import "testing"
import "github.com/SharedCode/parallels/database/repository"
import "os"

const defaultGroup = "1"

func TestUpsert(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	var config, _ = LoadConfiguration(dir + "/config.json")
	repo, e := NewRepository(config)
	if e != nil {
		t.Error(e)
	}
	rs := upsertData(repo)
	if rs.Error != nil {
		t.Error(rs.Error)
	}
}

func upsertData(repo repository.Repository) repository.Result {
	rs := repo.Set(*repository.NewKeyValue(defaultGroup, "K1", []byte("testV")))
	return rs
}

func TestRead(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	var config, _ = LoadConfiguration(dir + "/config.json")
	repo, e := NewRepository(config)
	if e != nil {
		t.Error(e)
	}

	// ensure we have data to read
	//upsertData(repo)

	r,rs := repo.Get(defaultGroup, "K1")
	if r != nil {
		if string(r[0].Value) != "testV" {
			t.FailNow()
		}
	} else {
		t.Error("K1 did not read from Store.")		
	}

	if rs.Error != nil {
		t.Error(rs.Error)
	}
}

func TestDelete(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	var config, _ = LoadConfiguration(dir + "/config.json")
	repo, e := NewRepository(config)
	if e != nil {
		t.Error(e)
	}
	// ensure we have data to read
	upsertData(repo)

	rs := repo.Remove(defaultGroup, "K1")
	if !rs.IsSuccessful() {
		t.Error(rs.Error)
	}
	gr,rs := repo.Get(defaultGroup, "K1")
	if !rs.IsSuccessful() {
		t.Error(rs.Error)
	}
	if gr != nil {
		t.Errorf("Expected K1 row to be deleted, but still found.")
	}
}

func TestNavigation(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	var config, _ = LoadConfiguration(dir + "/config.json")
	repoSet, e := NewRepositorySet(config)
	if e != nil {
		t.Error(e)
	}

	// ensure we have data to read
	rs := upsertData(repoSet.NavigableStore)
	if !rs.IsSuccessful() {
		t.Error(rs.Error)
	}

	// navigate to retrieve 1st "batch".
	r,rs := repoSet.NavigableStore.Navigate(defaultGroup, repository.Filter{})
	if rs.Error != nil {
		t.Error(rs.Error)
	}
	if r == nil {
		t.Error("Expected returned Result not found.")
		return
	}
	for _, kvp := range r {
		if kvp.Key != "K1" {
			t.Error("K1 key not found.")
		}
		if string(kvp.Value) != "testV" {
			t.Error("testV value not found.")
		}
	}
}

func TestNavigableDelete(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	var config, _ = LoadConfiguration(dir + "/config.json")
	repoSet, e := NewRepositorySet(config)
	if e != nil {
		t.Error(e)
	}
	repo := repoSet.NavigableStore
	// ensure we have data to read
	upsertData(repo)

	rs := repo.Remove(defaultGroup, "K1")
	if !rs.IsSuccessful() {
		t.Error(rs.Error)
	}
	r,rs := repo.Get(defaultGroup, "K1")
	if rs.Error != nil {
		t.Error(rs.Error)
	}
	if r != nil {
		t.Errorf("Expected K1 row to be deleted, but still found.")
	}
}
