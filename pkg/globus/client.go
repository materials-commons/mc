package globus

import (
	"crypto/tls"
	"fmt"

	"gopkg.in/resty.v1"

	"github.com/pkg/errors"
)

var ErrGlobusAuth = errors.New("globus auth")

type Client struct {
	token      string
	ccUser     string
	ccPassword string
}

const authURLBase = "https://auth.globus.org/v2"
const transferManagerURLBase = "https://transfer.api.globus.org/v0.10"

func CreateConfidentialClient(ccUser, ccPassword string) (*Client, error) {
	c := &Client{
		ccUser:     ccUser,
		ccPassword: ccPassword,
	}

	if err := c.Authenticate(); err != nil {
		return nil, errors.Wrapf(err, "Unable to authenticate to Globus")
	}

	return c, nil
}

// Authenticate will authentice with Globus and get a token to use in subsequent calls. The format of
// the response is:
//
// {"access_token":"AgxBpPJNQ98VY1Q6zk7gn43Y6rnzBywwJD2VzKlVdDjpQYDvV2u8Cleyd0DG1QlwXk2DM3jDdzjl2YfqGmqqghK134",
// "expires_in":172800,
// "resource_server":"transfer.api.globus.org",
// "token_type":"Bearer",
// "other_tokens":[],
// "scope":"urn:globus:auth:scope:transfer.api.globus.org:all"}
func (c *Client) Authenticate() error {
	var authResp struct {
		Token          string   `json:"access_token"`
		ExpiresIn      int      `json:"expires_in"`
		ResourceServer string   `json:"resource_server"`
		TokenType      string   `json:"token_type"`
		OtherTokens    []string `json:"other_tokens"`
		Scope          string   `json:"scope"`
	}
	resp, err := r().
		SetBasicAuth(c.ccUser, c.ccPassword).
		SetResult(&authResp).
		SetQueryParams(map[string]string{
			"grant_type": "client_credentials",
			"scope":      "urn:globus:auth:scope:transfer.api.globus.org:all",
		}).Post(fmt.Sprintf("%s/oauth2/token", authURLBase))

	if err := getAPIError(resp, err); err != nil {
		return err
	}

	c.token = authResp.Token
	return nil
}

func (c *Client) GetEndpointTaskList(endpointID string, filters map[string]string) (TaskList, error) {
	queryFilters := make(map[string]string)
	queryFilters["filter_endpoint"] = endpointID
	if filters != nil {
		for key, value := range filters {
			queryFilters[key] = value
		}
	}

	var taskList TaskList
	request := r().SetAuthToken(c.token).SetQueryParams(queryFilters).SetResult(&taskList)
	url := fmt.Sprintf("%s/endpoint_manager/task_list", transferManagerURLBase)

	resp, err := request.Get(url)
	err = getAPIError(resp, err)
	if err == ErrGlobusAuth {
		err = c.reauthAndRedoGet(request, url)
	}
	return taskList, err
}

func (c *Client) GetTaskSuccessfulTransfers(taskID string, marker int) (TransferItems, error) {
	var items TransferItems
	request := r().SetAuthToken(c.token).SetResult(&items)
	if marker != 0 {
		request = request.SetQueryParam("marker", fmt.Sprintf("%d", marker))
	}
	url := fmt.Sprintf("%s/endpoint_manager/task/%s/successful_transfers", transferManagerURLBase, taskID)

	resp, err := request.Get(url)
	err = getAPIError(resp, err)
	if err == ErrGlobusAuth {
		err = c.reauthAndRedoGet(request, url)
	}
	return items, err
}

func (c *Client) reauthAndRedoGet(request *resty.Request, url string) error {
	if err := c.Authenticate(); err != nil {
		return err
	}

	resp, err := request.Get(url)
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
	case resp.RawResponse.StatusCode == 401:
		return ErrGlobusAuth
	case resp.RawResponse.StatusCode > 299:
		return errors.New("unable to connect")
	default:
		return nil
	}
}
