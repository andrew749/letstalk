package api

type Credential struct {
	Id   uint   `json:"id" binding:"required"`
	Name string `json:"name" binding:"required"`
}

type AddUserCredentialRequest struct {
	Name string `json:"name" binding:"required"`
}

type AddUserCredentialResponse struct {
	CredentialId uint `json:"credentialId" binding:"required"`
}

type RemoveUserCredentialRequest struct {
	CredentialId uint `json:"credentialId" binding:"required"`
}

// If you provide a non-zero id, adds id, otherwise uses name if non-empty.
// Mainly for backwards compatibility.
type AddUserCredentialRequestRequest struct {
	CredentialId uint   `json:"credentialId"`
	Name         string `json:"name"`
}

type AddUserCredentialRequestResponse struct {
	CredentialId uint `json:"credentialId" binding:"required"`
}

type RemoveUserCredentialRequestRequest struct {
	CredentialId uint `json:"credentialId" binding:"required"`
}
