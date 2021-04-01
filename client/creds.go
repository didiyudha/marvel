package client

import (
	"crypto/md5"
	"encoding/hex"
)

type CredentialRequest struct {
	Host string
	PublicKey  string
	PrivateKey string
}

func (c *CredentialRequest) Hash(reqTimeStamps string) string {
	hasher := md5.New()
	hasher.Write([]byte(reqTimeStamps + c.PrivateKey + c.PublicKey))
	return hex.EncodeToString(hasher.Sum(nil))
}

func (c *CredentialRequest) APIKey() string {
	return c.PublicKey
}