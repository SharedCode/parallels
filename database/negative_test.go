package database

import "testing"
import "os"

func TestNullUpsert(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	var config, _ = LoadConfiguration(dir + "/config.json")
	repo, e := NewRepository(config)
	if e != nil {
		t.Error(e)
	}
	rs := repo.Set(nil...)
	if !rs.IsSuccessful() {
		t.Error(rs.Error)
	}
}

func TestNullGet(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	var config, _ = LoadConfiguration(dir + "/config.json")
	repo, e := NewRepository(config)
	if e != nil {
		t.Error(e)
	}
	_,rs := repo.Get("0", "")
	if !rs.IsSuccessful() {
		t.Error(rs.Error)
	}
}

func TestNullDelete(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	var config, _ = LoadConfiguration(dir + "/config.json")
	repo, e := NewRepository(config)
	if e != nil {
		t.Error(e)
	}
	rs := repo.Remove("0", "")
	if !rs.IsSuccessful() {
		t.Error(rs.Error)
	}
}
