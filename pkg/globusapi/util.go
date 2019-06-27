package globusapi

import (
	"os"

	"github.com/apex/log"
)

type APIParams struct {
	GlobusEndpointID string
	GlobusCCUser     string
	GlobusCCToken    string
}

func GetAPIParamsFromEnvFatal() APIParams {
	apiParams := APIParams{
		GlobusEndpointID: os.Getenv("MC_CONFIDENTIAL_CLIENT_ENDPOINT"),
		GlobusCCUser:     os.Getenv("MC_CONFIDENTIAL_CLIENT_USER"),
		GlobusCCToken:    os.Getenv("MC_CONFIDENTIAL_CLIENT_PW"),
	}

	if apiParams.GlobusEndpointID == "" {
		log.Fatalf("MC_CONFIDENTIAL_CLIENT_ENDPOINT env var is unset")
	}

	if apiParams.GlobusCCUser == "" {
		log.Fatalf("MC_CONFIDENTIAL_CLIENT_USER env var is unset")
	}

	if apiParams.GlobusCCToken == "" {
		log.Fatalf("MC_CONFIDENTIAL_CLIENT_PW env var is unset")
	}

	return apiParams
}
