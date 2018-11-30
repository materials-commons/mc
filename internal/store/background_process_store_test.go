package store_test

import (
	"fmt"
	"testing"

	"github.com/materials-commons/mc/internal/store"
	"github.com/materials-commons/mc/internal/store/model"
	"github.com/materials-commons/mc/internal/store/storengine"

	r "gopkg.in/gorethink/gorethink.v4"
)

func TestBackgroundStore_AddBackgroundProcess(t *testing.T) {
	session, _ := r.Connect(r.ConnectOpts{Database: "mctest", Address: "localhost:30815"})
	r.SetTags("r")
	se := storengine.NewBackgroundProcessRethinkdb(session)
	bgp := store.NewBackgroundProcessStore(se)

	bgpAdd := model.AddBackgroundProcessModel{}

	load, err := bgp.AddBackgroundProcess(bgpAdd)
	fmt.Println(err)
	fmt.Printf("created %#v\n", load)

}
