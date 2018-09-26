package globus

import (
	"fmt"

	"github.com/pkg/errors"
)

type ConfidentialClient struct {
	token      string
	ccUser     string
	ccPassword string
}

const authURLBase = "https://auth.globus.org/v2"
const transferManagerURLBase = "https://transfer.api.globus.org/v0.10"

func CreateConfidentialClient(ccUser, ccPassword string) (*ConfidentialClient, error) {
	c := &ConfidentialClient{
		ccUser:     ccUser,
		ccPassword: ccPassword,
	}

	if err := c.Authenticate(); err != nil {
		return nil, errors.Wrapf(err, "Unable to authenticate to Globus")
	}

	return c, nil
}

func (c *ConfidentialClient) Authenticate() error {
	resp, err := r().
		SetBasicAuth(c.ccUser, c.ccPassword).
		SetQueryParams(map[string]string{
			"grant_type": "client_credentials",
			"scope":      "urn:globus:auth:scope:transfer.api.globus.org:all",
		}).Post(fmt.Sprintf("%s/oauth2/token", authURLBase))

	return getAPIError(resp, err)
}
