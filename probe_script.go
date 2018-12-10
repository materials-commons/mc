package main

import (
	"fmt"
	"os"
	"time"

	"github.com/materials-commons/mc/internal/store"
	"github.com/materials-commons/mc/internal/store/model"
	"github.com/materials-commons/mc/internal/store/storengine"
	r "gopkg.in/gorethink/gorethink.v4"
)

func main() {
	fmt.Println("Starting...")

	testProjectName := "Test1"

	userID := "test@test.mc"
	fmt.Println(userID)

	session, _ := r.Connect(r.ConnectOpts{Database: "materialscommons", Address: "localhost:30815"})
	r.SetTags("r")

	projectsStoreEngine := storengine.NewProjectsRethinkdb(session)
	projectsStore := store.NewProjectsStore(projectsStoreEngine)
	projectID := getTestProject(projectsStore, userID, testProjectName)
	if projectID == "" {
		fmt.Printf("Can not find test project %s\n", testProjectName)
		fmt.Println("Exiting.")
		os.Exit(-1)
	}
	fmt.Printf("Found project %s, id = %s\n", testProjectName, projectID)

	guStoreEngin := storengine.NewGlobusUploadsRethinkdb(session)
	globusUploadsStore := store.NewGlobusUploadsStore(guStoreEngin)

	fmt.Printf("Waiting on a globusUpload for user %s in project%s (%s)/n", userID, testProjectName, projectID)
	bgTaskID := ""
	for {
		globusUploadsList := getGlobusUploadsList(globusUploadsStore, userID, projectID)
		if len(globusUploadsList) > 0 {
			fmt.Printf("Found %d GlobusUpload record(s)\n", len(globusUploadsList))
			for _, gul := range globusUploadsList {
				fmt.Printf("Found GlobusUpload with id = %s\n", gul.ID)
			}
			bgTaskID = globusUploadsList[0].ID
			break
		}
		time.Sleep(100)
	}

	fmt.Printf("Waiting on success status for %s\n", bgTaskID)

	bgpStoreEngine := storengine.NewBackgroundProcessRethinkdb(session)
	bgps := store.NewBackgroundProcessStore(bgpStoreEngine)

	getListModel := model.GetListBackgroundProcessModel{
		UserID:           userID,
		ProjectID:        projectID,
		BackgroundTaskID: bgTaskID,
	}

	var finished bool
	finished = false
	for {
		fmt.Println("status loop")
		time.Sleep(5 * time.Second)

		bgpList, err := bgps.GetListBackgroundProcess(getListModel)
		if err != nil {
			fmt.Printf("Unexpected error %+v\n", err)
			break
		}
		if len(bgpList) > 0 {
			backgroundProcess := bgpList[0]
			fmt.Println("First background process record...")
			fmt.Printf("   status = %s\n   message = %s\n", backgroundProcess.Status, backgroundProcess.Message)
			fmt.Printf("   success = %t\n", backgroundProcess.IsOk)
			finished = backgroundProcess.IsFinished
		}
		if finished {
			break
		}
	}

	fmt.Println("Done.")
}

func getTestProject(projectsStore *store.ProjectsStore, userID string, projectName string) string {
	var probe *model.ProjectCountModel

	projects, _ := projectsStore.GetProjectsForUser(userID)
	for _, project := range projects {
		if projectName == project.Name {
			probe = &project
		}
	}
	if probe != nil {
		return probe.ID
	}
	return ""
}

func getGlobusUploadsList(gus *store.GlobusUploadsStore, userID string, projectID string) []model.GlobusUploadSchema {
	var ret []model.GlobusUploadSchema
	ret, _ = gus.GetAllGlobusUploadsForUser(userID)
	return ret
}
