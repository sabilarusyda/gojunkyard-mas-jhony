package csurf

import (
	"crypto/sha1"
	"encoding/base64"
	"math/rand"
	"strings"
	"time"
)

var (
	Csrf_  Csrf
	random = rand.New(rand.NewSource(time.Now().UnixNano()))
)

type (
	// Csrf dependency
	Csrf struct {
		SaltLength   uint // 8 The string length of the salt
		SecretLength uint // 18 The byte length of the secret key
	}
)

// Tokens ...
// Token generation/verification class.
func (c *Csrf) Tokens(opts *Options) *Csrf {
	if opts.SaltLength == 0 {
		opts.SaltLength = 8
	}
	c.SaltLength = opts.SaltLength

	if opts.SecretLength == 0 {
		opts.SecretLength = 18
	}
	c.SecretLength = opts.SecretLength

	return c
}

// Create ...
// Create a new CSRF token.
func (c *Csrf) Create(secret string) string {
	return c.Tokenize(secret, c.Rndm(c.SaltLength))
}

// Verify ...
// Verify CSRF token
func (c *Csrf) Verify(secret, token string) bool {
	idx := strings.IndexByte(token, '-')
	if idx < 0 {
		return false
	}

	expected := c.Tokenize(secret, token[:idx])
	return expected == token
}

// Rndm ...
// Generate random string
// base62 = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789'
// base36 = 'abcdefghijklmnopqrstuvwxyz0123456789'
// base10 = '0123456789'
func (c *Csrf) Rndm(length uint) string {
	var (
		salt  string
		chars = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")
	)

	if length == 0 {
		length = 10
	}

	for i := 0; i < int(length); i++ {
		salt = salt + string(chars[random.Intn(len(chars))])
	}

	return salt
}

// Tokenize a secret and salt.
func (c *Csrf) Tokenize(secret string, salt string) string {
	return salt + "-" + c.Hash(salt+"-"+secret)
}

// Hash a string with SHA1, returning url-safe base64
func (c *Csrf) Hash(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	bs := h.Sum(nil)
	s = base64.StdEncoding.EncodeToString(bs)

	s = strings.Replace(s, "+", "-", -1)
	s = strings.Replace(s, "/", "_", -1)
	s = strings.Replace(s, "=", "", -1)

	return s
}
