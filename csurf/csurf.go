package csurf

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

type (
	// Options ...
	Options struct {
		Cookie        Cookie
		SessionKey    string
		Value         string
		IgnoreMethods []string
		SaltLength    uint // 8 The string length of the salt
		SecretLength  uint // 18 The byte length of the secret key
	}

	// Cookie ...
	Cookie struct {
		Key    string
		Path   string
		Signed bool
	}
)

// Csurf is CSRF protection middleware.
// This middleware adds a `req csrfToken()` function to make a token
// which should be added to requests which mutate
// state, within a hidden form field, query-string etc. This
// token is validated against the visitor's session.
func (opts *Options) Csurf(r *http.Request) {
	// get cookie options
	cookie := opts.GetCookieOptions()

	// get session options
	if opts.SessionKey == "" {
		opts.SessionKey = "session"
	}
	sessionKey := opts.SessionKey

	// get value getter
	if opts.Value == "" {
		opts.Value = opts.DefaultValue(r)
	}
	value := opts.Value

	// token repo
	token := Csrf_.Tokens(opts)

	// ignored methods
	if len(opts.IgnoreMethods) == 0 {
		opts.IgnoreMethods = []string{"GET", "HEAD", "OPTIONS"}
	}
	ignoreMethods := opts.IgnoreMethods

	// generate lookup
	ignoreMethod := opts.GetIgnoredMethods(ignoreMethods)

	fmt.Println(cookie, sessionKey, value, token, ignoreMethods, ignoreMethod)
}

// DefaultValue function, checking the req body
// and req query for the CSRF token.
func (opts *Options) DefaultValue(r *http.Request) string {
	if value := r.URL.Query().Get("_csrf"); value != "" {
		return value
	} else if value = r.FormValue("_csrf"); value != "" {
		return value
	} else if value = r.Header.Get("csrf-token"); value != "" {
		return value
	} else if value = r.Header.Get("xsrf-token"); value != "" {
		return value
	} else if value = r.Header.Get("x-csrf-token"); value != "" {
		return value
	} else if value = r.Header.Get("x-xsrf-token"); value != "" {
		return value
	}

	return ""
}

// GetCookieOptions is method to set default value for path and key of cookie
// set key to `_csrf` as default value if empty
// set path to `/` as default value if empty
func (opts *Options) GetCookieOptions() Cookie {
	if opts.Cookie.Key == "" {
		opts.Cookie.Key = "_csrf"
	} else {
		opts.Cookie.Signed = true
	}

	if opts.Cookie.Path == "" {
		opts.Cookie.Path = "/"
	}

	return opts.Cookie
}

// GetIgnoredMethods ...
// Get a lookup of ignored methods.
func (opts *Options) GetIgnoredMethods(methods []string) map[string]bool {
	var obj map[string]bool

	for _, v := range methods {
		method := strings.ToUpper(v)
		obj[method] = true
	}

	return obj
}

// GetSecret ...
// Get the token secret from the request.
func (opts *Options) GetSecret(r *http.Request, sessionKey string) (string, error) {
	// get the bag & key
	bag := opts.GetSecretBag(r, sessionKey)
	if opts.Cookie.Key == "" {
		opts.Cookie.Key = "csrfSecret"
	}
	key := opts.Cookie.Key

	if len(bag) == 0 {
		return "", errors.New("Error")
	}

	// return secret from bag
	return bag[key], nil
}

// GetSecretBag ...
// Sample return { _csrf: '2gEGvMDvm2nY1zF8B8zeNc4V' }
func (opts *Options) GetSecretBag(r *http.Request, sessionKey string) map[string]string {
	if opts.Cookie.Signed {
		cookie, err := r.Cookie("signedCookies")
		if err != nil {
			return map[string]string{}
		}

		return map[string]string{opts.Cookie.Key: cookie.Value}
	} else if !opts.Cookie.Signed {
		cookie, err := r.Cookie("cookies")
		if err != nil {
			return map[string]string{}
		}

		return map[string]string{opts.Cookie.Key: cookie.Value}
	}

	// get secret from session
	cookie, err := r.Cookie(sessionKey)
	if err != nil {
		return map[string]string{}
	}

	return map[string]string{opts.Cookie.Key: cookie.Value}
}

// SetCookie ...
// Set a cookie on the HTTP response.
func (opts *Options) SetCookie(w http.ResponseWriter, val string, options map[string]string) {
	expire := time.Now().AddDate(0, 1, 0)
	if options["path"] == "" {
		options["path"] = "/"
	}

	cookie := http.Cookie{
		Name:    opts.Cookie.Key,
		Value:   val,
		Expires: expire,
		Path:    options["path"],
	}
	http.SetCookie(w, &cookie)
}

// SetSecret ...
// Set the token secret on the request.
func (opts *Options) SetSecret(w http.ResponseWriter) string {
	var (
		secret string
		chars  = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")
		length = opts.SecretLength
	)

	rand.Seed(time.Now().UnixNano())

	if length == 0 {
		length = 24
	}

	for i := 0; i < int(length); i++ {
		secret = secret + string(chars[rand.Intn(len(chars))])
	}

	opts.SetCookie(w, secret, map[string]string{"path": opts.Cookie.Path})

	return secret
}

// VerifyConfiguration ...
// Verify the configuration against the request.
func (opts *Options) VerifyConfiguration(r *http.Request, sessionKey string) bool {
	if sb := opts.GetSecretBag(r, sessionKey); len(sb) == 0 {
		return false
	}

	if opts.Cookie.Signed == true && r.URL.Query().Get("secret") == "" {
		return false
	}

	return true
}
