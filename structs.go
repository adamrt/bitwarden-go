package main

type tokenPayload struct {
	// GrantType determines the code path. [refresh_token,password]
	GrantType string `schema:"grant_type"`

	// RefreshToken is used for grant_type == refresh_token
	RefreshToken string `schema:"refresh_token"`

	// These are used for grant_type == password
	ClientID         string `schema:"client_id"`
	Username         string `schema:"username"`
	Password         string `schema:"password"`
	Scope            string `schema:"scope"`
	DeviceIdentifier string `schema:"deviceidentifier"`
	DeviceName       string `schema:"devicename"`
	DeviceType       string `schema:"devicetype"`
	DevicePushToken  string `schema:"devicepushtoken"`
}

type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	Key          string `json:"key"`
}

type cipherPayload struct {
	Type             string `schema:"type"`
	Name             string `schema:"name"`
	FolderUUID       string `schema:"folderid"`
	OrganizationUUID string `schema:"organizationid"`
	Favorite         bool   `schema:"favorite"`
}

type cipherDataPayload struct {
	Name string
}

type registerJSON struct {
	MasterPasswordHash string `json:"masterpasswordhash"`
	MasterPasswordHint string `json:"masterpasswordhint"`
	Email              string `json:"email"`
	Key                string `json:"key"`
	Password           string `json:"password"`
	Scope              string `json:"scope"`
	DeviceIdentifier   string `json:"deviceidentifier"`
	DeviceName         string `json:"devicename"`
	DeviceType         string `json:"devicetype"`
	DevicePushToken    string `json:"devicepushtoken"`
}

type syncResponse struct {
	Profile User
	Folders []Folder
	Ciphers []Cipher
	Object  string
	Domains struct {
		EquivalentDomains       *[]string
		GlobalEquivalentDomains []string
		Object                  string
	}
}
