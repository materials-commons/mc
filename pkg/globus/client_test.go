package globus_test

import (
	"os"
	"testing"

	"github.com/materials-commons/mc/pkg/globus"
	"github.com/materials-commons/mc/pkg/tutils/assert"
)

func TestGetIdentities(t *testing.T) {
	globusCCUser := os.Getenv("MC_CONFIDENTIAL_CLIENT_USER")
	globusCCToken := os.Getenv("MC_CONFIDENTIAL_CLIENT_PW")

	client, err := globus.CreateConfidentialClient(globusCCUser, globusCCToken)
	assert.Okf(t, err, "Unable to create confidential client: %s", err)

	identities, err := client.GetIdentities([]string{"glenn.tarcea@gmail.com"})
	assert.Okf(t, err, "Unable to get identities: %s", err)
	assert.Truef(t, len(identities.Identities) == 1, "Wrong identities length %d", len(identities.Identities))
}
