package store_test

import (
	"fmt"
	"testing"

	"github.com/materials-commons/mc/internal/store"
	"github.com/materials-commons/mc/internal/store/model"
	"github.com/materials-commons/mc/internal/store/storengine"

	r "gopkg.in/gorethink/gorethink.v4"
)

func TestFileLoadsStore_AddFileLoad(t *testing.T) {
	if true {
		t.Skip()
	}
	session, _ := r.Connect(r.ConnectOpts{Database: "mctest", Address: "localhost:30815"})
	r.SetTags("r")
	se := storengine.NewFileLoadsRethinkdb(session)
	floads := store.NewFileLoadsStore(se)

	flAdd := model.AddFileLoadModel{
		ProjectID: "abc123",
		Path:      "/tmp",
		Owner:     "test@test.com",
	}

	load, err := floads.AddFileLoad(flAdd)
	fmt.Println(err)
	fmt.Printf("created %#v\n", load)

}
