package globus

import (
	"crypto/tls"

	"errors"

	"gopkg.in/resty.v1"
)

func GetToken(userID, userPW string) error {
	resp, err := r().
		SetBasicAuth(userID, userPW).
		SetQueryParams(map[string]string{
			"grant_type": "client_credentials",
			"scope":      "urn:globus:auth:scope:transfer.api.globus.org:all",
		}).
		Post("https://auth.globus.org/v2/oauth2/token")
	return getAPIError(resp, err)
}

var tlsConfig = tls.Config{InsecureSkipVerify: true}

// r is similar to resty.R() except that it sets the TLS configuration
func r() *resty.Request {
	return resty.SetTLSClientConfig(&tlsConfig).R()
}

func getAPIError(resp *resty.Response, err error) error {
	switch {
	case err != nil:
		return err
	case resp.RawResponse.StatusCode > 299:
		return errors.New("unable to connect")
	default:
		return nil
	}
}
