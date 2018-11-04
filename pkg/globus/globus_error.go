package globus

import (
	"encoding/json"
	"fmt"

	"github.com/materials-commons/mc/pkg/mc"

	"github.com/pkg/errors"
	"gopkg.in/resty.v1"
)

// ErrorResponse describes the JSON that Globus responds with when there is an error in an API call
type ErrorResponse struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	RequestID string `json:"request_id"`
	Resource  string `json:"resource"`
}

func ToErrorFromResponse(resp *resty.Response) error {
	var errorResponse ErrorResponse
	if err := json.Unmarshal(resp.Body(), &errorResponse); err != nil {
		return errors.WithMessage(mc.ErrGlobusAPI, fmt.Sprintf("(HTTP Status: %d)- unable to parse json error response: %s", resp.RawResponse.StatusCode, err))
	}

	return errors.WithMessage(mc.ErrGlobusAPI, fmt.Sprintf("(HTTP Status: %d)- %s: %s", resp.RawResponse.StatusCode, errorResponse.Code, errorResponse.Message))
}
