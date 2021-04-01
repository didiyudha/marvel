package client

import (
	"github.com/bxcodec/faker/v3"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
	"time"
)

func TestHash(t *testing.T) {
	ts := strconv.FormatInt(time.Now().UTC().UnixNano(), 16)
	c := &CredentialRequest{
		Host:       faker.IPv4(),
		PublicKey:  faker.Username(),
		PrivateKey: faker.Email(),
	}
	hash := c.Hash(ts)
	assert.False(t, hash == "")
}

func TestAPIKey(t *testing.T) {
	c := &CredentialRequest{
		Host:       faker.IPv4(),
		PublicKey:  faker.Username(),
		PrivateKey: faker.Email(),
	}
	apiKey := c.APIKey()
	assert.Equal(t, c.PublicKey, apiKey)
}