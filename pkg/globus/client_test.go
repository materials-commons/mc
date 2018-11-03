package globus_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/materials-commons/mc/pkg/globus"
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

	rule := globus.EndpointACLRule{
		EndpointID:  testEndpointID,
		Path:        "",
		IdentityID:  identities.Identities[0].ID,
		Permissions: "rw",
	}
	aclRes, err := client.AddEndpointACLRule(rule)
	assert.Okf(t, err, "Unable to set ACL rule: %s", err)
	fmt.Printf("%#v\n", aclRes)

	delRes, err := client.DeleteEndpointACLRule(testEndpointID, aclRes.AccessID)
	assert.Okf(t, err, "Unable to delete ACL rule: %s", err)
	fmt.Printf("#%v\n", delRes)
}

func createClient(t *testing.T) *globus.Client {
	globusCCUser := os.Getenv("MC_CONFIDENTIAL_CLIENT_USER")
	globusCCToken := os.Getenv("MC_CONFIDENTIAL_CLIENT_PW")
	testEndpointID = os.Getenv("MC_CONFIDENTIAL_CLIENT_ENDPOINT")

	if globusCCUser != "" && globusCCToken != "" && testEndpointID != "" {
		client, err := globus.CreateConfidentialClient(globusCCUser, globusCCToken)
		assert.Okf(t, err, "Unable to create confidential client: %s", err)
		return client
	} else {
		t.Skipf("One or more of MC_CONFIDENTIAL_CLIENT_USER, MC_CONFIDENTIAL_CLIENT_PW, MC_CONFIDENTIAL_CLIENT_ENDPOINT not set")
		return nil
	}
}
