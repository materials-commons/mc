package globusapi

type Identities struct {
	Identities []Identity `json:"identities"`
}

type IncludedProviders struct {
	IdentityProviders []IdentityProvider `json:"identity_providers"`
	Identities        []Identity         `json:"identities"`
}

type IdentityProvider struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Identity struct {
	Username         string `json:"username"`
	Status           string `json:"status"`
	Name             string `json:"name"`
	ID               string `json:"id"`
	IdentityProvider string `json:"identity_provider"`
	Organizaion      string `json:"organizaton"`
	EMail            string `json:"email"`
}
