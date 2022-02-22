package pubsub

import (
	b64 "encoding/base64"
	"fmt"
	"net/http"
	"time"

	request "devcode.xeemore.com/systech/gojunkyard/http/httprequest"

	jwt "github.com/dgrijalva/jwt-go"
)

// Options for the pubsub.
type Options struct {
	Sub      string
	Issuer   string
	Secret   string
	Endpoint string
}

// Client object
type Client struct {
	Options Options
}

// NewPubsubClient create new pubsub client
func NewPubsubClient(options *Options) *Client {
	return &Client{
		Options: *options,
	}
}

//CreateAuthorizationToken generate authorization token for pubsub
func (c *Client) CreateAuthorizationToken() (string, error) {
	mySigningKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(c.Options.Secret))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.StandardClaims{
		Subject:   c.Options.Sub,
		Issuer:    c.Options.Issuer,
		Audience:  "https://pubsub.googleapis.com/google.pubsub.v1.Publisher",
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Minute * 65).Unix(),
	})

	ss, err := token.SignedString(mySigningKey)
	if err != nil {
		return "", err
	}

	return ss, nil
}

//SendIntoPubSub send object into pubsub
func (c *Client) SendIntoPubSub(token string, jsonObject []byte) (int, error) {
	type Message struct {
		Data string `json:"data"`
	}

	type ReqBody struct {
		Messages []Message `json:"messages"`
	}

	reqBody := ReqBody{
		Messages: []Message{
			Message{
				Data: b64.StdEncoding.EncodeToString(jsonObject),
			},
		},
	}

	hh := http.Header{}
	hh.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	hh.Set("Content-Type", "application/json")

	req := request.NewHTTPRequest()
	req.Header = hh
	req.Method = request.MethodPostJSON
	req.Body = reqBody
	req.URL = c.Options.Endpoint

	response, body, err := req.Send()
	if err != nil {
		return 500, err
	} else if response.StatusCode != 200 {
		return response.StatusCode, fmt.Errorf(string(body))
	}

	return 200, nil
}
