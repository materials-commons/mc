package globus_test

import (
    "os"
	"crypto/tls"
	"errors"
	"fmt"
	"testing"

	"gopkg.in/resty.v1"
)

func TestGetToken(t *testing.T) {
	if false {
		userID := os.Getenv("FOO")
		userPW := os.Getenv("FOO")
		resp, err := r().
			SetBasicAuth(userID, userPW).
			SetQueryParams(map[string]string{
				"grant_type": "client_credentials",
				"scope":      "urn:globus:auth:scope:transfer.api.globus.org:all",
			}).
			Post("https://auth.globus.org/v2/oauth2/token")
		fmt.Println("err", err)
		fmt.Printf("%+v\n", resp)
	}
}

func TestGetTaskList(t *testing.T) {
	if true {
		token := os.Getenv("FOO")
		resp, err := r().
			SetAuthToken(token).
			SetQueryParams(map[string]string{
				"filter_endpoint": os.Getenv("FOO"),
				"limit":           "1000",
			}).
			Get("https://transfer.api.globus.org/v0.10/endpoint_manager/task_list")
		fmt.Println("err", err)
		fmt.Printf("%+v\n", resp)
	}
}

func GetIdentities(t *testing.T) {
	//SetQueryParams(map[string]string{
	//	"identities": os.Getenv("FOO"),
	//	"provision":  "false",
	//}).
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
