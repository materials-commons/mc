package main

import (
	"fmt"
	"os"

	"github.com/materials-commons/mc/internal/store"
	"github.com/materials-commons/mc/internal/store/model"
	"github.com/materials-commons/mc/internal/store/storengine"
	r "gopkg.in/gorethink/gorethink.v4"
)

func main() {
	fmt.Println("Starting...")

	session, _ := r.Connect(r.ConnectOpts{Database: "mctest", Address: "localhost:30815"})
	r.SetTags("r")
	storeEngine := storengine.NewBackgroundProcessRethinkdb(session)
	bgps := store.NewBackgroundProcessStore(storeEngine)

	cleanupBackgroundProcessEngine(storeEngine)

	abgpModel := model.AddBackgroundProcessModel{
		UserID:                "bogues.user@mc.org",
		ProjectID:             "ProjectId",
		BackgroundProcessID:   "BGProcessId",
		BackgroundProcessType: "bgp-type",
		Status:                "status",
		Message:               "message",
	}

	getList := model.GetListBackgroundProcessModel{
		UserID:              "bogues.user@mc.org",
		ProjectID:           "ProjectId",
		BackgroundProcessID: "BGProcessId",
	}

	_, err := bgps.AddBackgroundProcess(abgpModel)
	if err != nil {
		fmt.Println("Error - All 1")
		os.Exit(-1)
	}
	_, err = bgps.AddBackgroundProcess(abgpModel)
	if err != nil {
		fmt.Println("Error - All 2")
		os.Exit(-1)
	}

	bgpList, err := bgps.GetListBackgroundProcess(getList)
	if err != nil {
		fmt.Println("Error - Get List")
		os.Exit(-1)
	}

	fmt.Printf("%v \n", len(bgpList))

	fmt.Println("Done.")
}

func cleanupBackgroundProcessEngine(e storengine.BackgroundProcessStoreEngine) {
	if re, ok := e.(*storengine.BackgroundProcessRethinkdb); ok {
		session := re.Session
		_, _ = r.Table("background_process").Delete().RunWrite(session)
	}
}
