package globusapi_test

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/materials-commons/mc/pkg/globusapi"
	"github.com/materials-commons/mc/pkg/tutils/assert"
)

var testEndpointID string

func TestClient_GetIdentities(t *testing.T) {
	client := createClient(t)
	identities, err := client.GetIdentities([]string{"glenn.tarcea@gmail.com"})
	assert.Okf(t, err, "Unable to get identities: %s", err)
	assert.Truef(t, len(identities.Identities) == 1, "Wrong identities length %d", len(identities.Identities))
}

func TestACLs(t *testing.T) {
	client := createClient(t)
	identities, err := client.GetIdentities([]string{"glenn.tarcea@gmail.com"})
	assert.Okf(t, err, "Unable to get identities: %s", err)
	userGlobusIdentity := identities.Identities[0].ID

	tests := []struct {
		identity   string
		path       string
		shouldFail bool
		deleteACL  bool
		name       string
	}{
		{identity: userGlobusIdentity, path: "/~/globus-staging/", shouldFail: false, deleteACL: false, name: "Add New ACL"},
		{identity: userGlobusIdentity, path: "/~/globus-staging/", shouldFail: false, deleteACL: true, name: "Add Existing ACL"},
		{identity: userGlobusIdentity, path: "/~/globus-staging", shouldFail: true, deleteACL: false, name: "Bad Path"},
	}

	for _, test := range tests {
		rule := globusapi.EndpointACLRule{
			PrincipalType: "identity",
			EndpointID:    testEndpointID,
			Path:          test.path,
			IdentityID:    test.identity,
			Permissions:   "rw",
		}

		aclRes, err := client.AddEndpointACLRule(rule)
		if !test.shouldFail {
			assert.Okf(t, err, "Unable to set ACL rule: %s - %#v", err, client.GetGlobusErrorResponse())
			if test.deleteACL {
				_, err := client.DeleteEndpointACLRule(testEndpointID, aclRes.AccessID)
				assert.Okf(t, err, "Unable to delete ACL rule: %s - %#v", err, client.GetGlobusErrorResponse())
			}
		} else {
			assert.Errorf(t, err, "Test should have failed")
		}
	}
}

func TestGetTasks(t *testing.T) {
	client := createClient(t)
	lastWeek := time.Now().AddDate(0, 0, -10).Format("2006-01-02")
	fmt.Println("lastWeek", lastWeek)
	tasks, err := client.GetEndpointTaskList(testEndpointID, map[string]string{
		"filter_completion_time": lastWeek,
		"filter_status":          "SUCCEEDED",
	})
	fmt.Println("GetEndpointTaskList err", err)
	fmt.Printf("   tasks: %#v\n", tasks)
	for _, task := range tasks.Tasks {
		transfers, err := client.GetTaskSuccessfulTransfers(task.TaskID, 0)
		fmt.Println("  GetTaskSuccessfulTransfers err", err)
		fmt.Printf("    transfers: %#v\n", transfers)
		transferItem := transfers.Transfers[0]
		pieces := strings.Split(transferItem.DestinationPath, "/")
		fmt.Println(len(pieces))
		fmt.Printf("pieces[0] = '%s'\n", pieces[0])
		fmt.Println("id =", pieces[2])
	}
}

func createClient(t *testing.T) *globusapi.Client {
	globusCCUser := os.Getenv("MC_CONFIDENTIAL_CLIENT_USER")
	globusCCToken := os.Getenv("MC_CONFIDENTIAL_CLIENT_PW")
	testEndpointID = os.Getenv("MC_CONFIDENTIAL_CLIENT_ENDPOINT")

	if globusCCUser != "" && globusCCToken != "" && testEndpointID != "" {
		client, err := globusapi.CreateConfidentialClient(globusCCUser, globusCCToken)
		assert.Okf(t, err, "Unable to create confidential client: %s", err)
		return client
	} else {
		t.Skipf("One or more of MC_CONFIDENTIAL_CLIENT_USER, MC_CONFIDENTIAL_CLIENT_PW, MC_CONFIDENTIAL_CLIENT_ENDPOINT not set")
		return nil
	}
}
