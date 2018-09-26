package globus_test

import (
	"crypto/tls"
	"errors"
	"fmt"
	"testing"

	"gopkg.in/resty.v1"
)

func TestGetToken(t *testing.T) {
	userID := "d8eba3ad-fdd9-468d-95a1-d5c4ff91de3f"
	userPW := "K3DWTaRt2MT9QzYA+Ttxs++9kea0cUO219wzaElQqP4="
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

func TestGetTaskList(t *testing.T) {
	token := "Ag3dv92pXXnbjmjlyG4ar22Ex0ez69m8Kne5DyrDGaQOQWzrlKCVCYddKKnDvbkg73OrYx3YzEVDVniXym3NVCdD7r"
	resp, err := r().
		SetAuthToken(token).
		SetQueryParams(map[string]string{
			"filter_endpoint":        "4e9d8294-bdcd-11e8-8c1e-0a1d4c5c824a",
			"filter_completion_time": "2018-09-23",
			"filter_status":          "SUCCEEDED",
			"limit":                  "1000",
		}).
		Get("https://transfer.api.globus.org/v0.10/endpoint_manager/task_list")
	fmt.Println("err", err)
	fmt.Printf("%+v\n", resp)
}

func GetIdentities(t *testing.T) {
	//SetQueryParams(map[string]string{
	//	"identities": "glenn.tarcea@gmail.com",
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
