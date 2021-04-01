package client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"strconv"
	"time"
)

func NewMarvelClient(cred CredentialRequest, restyClient *resty.Client) MarvelClient {
	c := marvelClientImpl{
		Client: restyClient,
		Creds:  &cred,
	}
	return &c
}

type MarvelClient interface {
	Characters(ctx context.Context, timestamps time.Time, limit, offset int) (CharacterResponse, error)
	CharacterDetail(ctx context.Context, timestamps time.Time, id int) (CharacterResponse, error)
}

type marvelClientImpl struct {
	Client *resty.Client
	Creds  *CredentialRequest
}

func (c *marvelClientImpl) CharactersFullURL(timestamps time.Time, limit, offset int) string {

	requestTimeStamp := strconv.FormatInt(timestamps.UTC().UnixNano(), 16)
	baseURL := fmt.Sprintf("http://%s/v1/public", c.Creds.Host)

	fullUrl := baseURL + fmt.Sprintf(
		`/characters?apikey=%s&ts=%s&hash=%s&limit=%d&offset=%d`,
		c.Creds.APIKey(),
		requestTimeStamp,
		c.Creds.Hash(requestTimeStamp),
		limit,
		offset,
	)

	return fullUrl
}

func (c *marvelClientImpl) Characters(ctx context.Context, timestamps time.Time, limit, offset int) (CharacterResponse, error) {

	fullUrl := c.CharactersFullURL(timestamps, limit, offset)

	res, err := c.Client.R().SetContext(ctx).Get(fullUrl)
	if err != nil {
		return CharacterResponse{}, err
	}

	var characterResponse CharacterResponse

	if err := json.Unmarshal(res.Body(), &characterResponse); err != nil {
		return CharacterResponse{}, errors.Wrap(err, "decode response body")
	}

	return characterResponse, nil
}

func (c *marvelClientImpl)  CharacterDetailFullURL(timestamps time.Time, id int) string {
	requestTimeStamp := strconv.FormatInt(timestamps.UTC().UnixNano(), 16)
	baseURL := fmt.Sprintf("http://%s/v1/public", c.Creds.Host)

	fullUrl := baseURL + fmt.Sprintf(
		`/characters/%d?apikey=%s&ts=%s&hash=%s`,
		id,
		c.Creds.APIKey(),
		requestTimeStamp,
		c.Creds.Hash(requestTimeStamp),
	)

	return fullUrl
}

func (c *marvelClientImpl) CharacterDetail(ctx context.Context, timestamps time.Time, id int) (CharacterResponse, error) {

	fullUrl := c.CharacterDetailFullURL(timestamps, id)

	res, err := c.Client.R().SetContext(ctx).Get(fullUrl)
	if err != nil {
		return CharacterResponse{}, err
	}

	var characterResponse CharacterResponse

	if err := json.Unmarshal(res.Body(), &characterResponse); err != nil {
		return CharacterResponse{}, errors.Wrap(err, "decode response body")
	}

	return characterResponse, nil
}
