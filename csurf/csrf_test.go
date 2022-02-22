package csurf

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCsrfCreateValid(t *testing.T) {
	var (
		secret = "qwertyuiop1234567890ASDF"
		salt   string
		aToken string // actual token
		aCsrf  string // actual csrf
		eToken string // expected token
		eCsrf  string // expected csrf

	)

	t.Log("Running test case `Create CSRF`")

	Csrf_.SaltLength = 8
	aToken = Csrf_.Create(secret)
	assert.IsType(t, aToken, "string")

	salt = aToken[:8]
	aCsrf = aToken[9:]

	eToken = salt + "-" + Csrf_.Hash(salt+"-"+secret)
	eCsrf = eToken[9:]

	assert.Equal(t, eToken, aToken)
	assert.Equal(t, eCsrf, aCsrf)

	ok := Csrf_.Verify(secret, aToken)
	assert.True(t, ok)
}

func TestCsrfRandomValid(t *testing.T) {
	t.Log("Running test case `Not equal random salt 10 chars`")
	assert.NotEqual(t, Csrf_.Rndm(10), Csrf_.Rndm(10))
}
