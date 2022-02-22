package httprequest

import (
	// "encoding/json"
	// "fmt"
	"encoding/xml"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/url"
	"testing"
)

type emptyObject struct{}
type normalObject struct {
	Key1 string `json:"key_1"`
}

func TestGetRequest(t *testing.T) {
	tests := []struct {
		name        string
		url         string
		params      url.Values
		header      http.Header
		want_status int
	}{
		{
			name:        "TEST GET",
			url:         "http://httpbin.org/get?show_env=1",
			params:      url.Values{"show_env": {"0"}},
			header:      http.Header{"Cache-Control": {"no-cache", "no-store", "must-revalidate"}},
			want_status: 200,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			header := tt.header
			params := tt.params
			httprequest := NewHTTPRequest()
			httprequest.Header = header
			httprequest.Method = MethodGet
			httprequest.URL = tt.url
			httprequest.Params = params

			response, _, _ := httprequest.Send()
			assert.Equal(t, tt.want_status, response.StatusCode)
		})
	}
}

func TestPostRequest(t *testing.T) {
	jsonex := struct {
		UserID   int64  `json:"user_id"`
		Username string `json:"username"`
	}{int64(999999), "robobot"}

	xmlex := struct {
		XMLName  xml.Name `xml:"user"`
		UserID   int64    `xml:"user_id"`
		Username string   `xml:"username"`
	}{UserID: int64(999999), Username: "robobot"}

	tests := []struct {
		name        string
		url         string
		params      url.Values
		header      http.Header
		meth        method
		body_data   interface{}
		want_status int
	}{
		{
			name:        "TEST STANDARD POST",
			url:         "http://httpbin.org/post",
			params:      url.Values{"show_env": {"0"}},
			header:      http.Header{"Cache-Control": {"no-cache", "no-store", "must-revalidate"}},
			meth:        MethodPost,
			body_data:   nil,
			want_status: 200,
		},
		{
			name:        "TEST BYTE POST",
			url:         "http://httpbin.org/post",
			params:      url.Values{"show_env": {"0"}},
			header:      http.Header{"Cache-Control": {"no-cache", "no-store", "must-revalidate"}},
			meth:        MethodPostRaw,
			body_data:   []byte("this is byte"),
			want_status: 200,
		},
		{
			name:        "TEST STRING POST",
			url:         "http://httpbin.org/post",
			params:      url.Values{"show_env": {"0"}},
			header:      http.Header{"Cache-Control": {"no-cache", "no-store", "must-revalidate"}},
			meth:        MethodPostRaw,
			body_data:   "this is string",
			want_status: 200,
		},
		{
			name:        "TEST JSON POST",
			url:         "http://httpbin.org/post",
			params:      url.Values{"show_env": {"0"}},
			header:      http.Header{"Cache-Control": {"no-cache", "no-store", "must-revalidate"}},
			meth:        MethodPostJSON,
			body_data:   jsonex,
			want_status: 200,
		},
		{
			name:        "TEST XML POST",
			url:         "http://httpbin.org/post",
			params:      url.Values{"show_env": {"0"}},
			header:      http.Header{"Cache-Control": {"no-cache", "no-store", "must-revalidate"}, "Content-Type": {"application/xml"}},
			meth:        MethodPostXML,
			body_data:   xmlex,
			want_status: 200,
		},
		{
			name:        "TEST STANDARD PUT",
			url:         "http://httpbin.org/put",
			params:      url.Values{"show_env": {"0"}},
			header:      http.Header{"Cache-Control": {"no-cache", "no-store", "must-revalidate"}},
			meth:        MethodPut,
			body_data:   nil,
			want_status: 200,
		},
		{
			name:        "TEST BYTE PUT",
			url:         "http://httpbin.org/put",
			params:      url.Values{"show_env": {"0"}},
			header:      http.Header{"Cache-Control": {"no-cache", "no-store", "must-revalidate"}},
			meth:        MethodPutRaw,
			body_data:   []byte("this is byte"),
			want_status: 200,
		},
		{
			name:        "TEST STRING PUT",
			url:         "http://httpbin.org/put",
			params:      url.Values{"show_env": {"0"}},
			header:      http.Header{"Cache-Control": {"no-cache", "no-store", "must-revalidate"}},
			meth:        MethodPutRaw,
			body_data:   "this is string",
			want_status: 200,
		},
		{
			name:        "TEST JSON PUT",
			url:         "http://httpbin.org/put",
			params:      url.Values{"show_env": {"0"}},
			header:      http.Header{"Cache-Control": {"no-cache", "no-store", "must-revalidate"}},
			meth:        MethodPutJSON,
			body_data:   jsonex,
			want_status: 200,
		},
		{
			name:        "TEST XML PUT",
			url:         "http://httpbin.org/put",
			params:      url.Values{"show_env": {"0"}},
			header:      http.Header{"Cache-Control": {"no-cache", "no-store", "must-revalidate"}, "Content-Type": {"application/xml"}},
			meth:        MethodPutXML,
			body_data:   xmlex,
			want_status: 200,
		},
		{
			name:        "TEST STANDARD PATCH",
			url:         "http://httpbin.org/patch",
			params:      url.Values{"show_env": {"0"}},
			header:      http.Header{"Cache-Control": {"no-cache", "no-store", "must-revalidate"}},
			meth:        MethodPatch,
			body_data:   nil,
			want_status: 200,
		},
		{
			name:        "TEST BYTE PATCH",
			url:         "http://httpbin.org/patch",
			params:      url.Values{"show_env": {"0"}},
			header:      http.Header{"Cache-Control": {"no-cache", "no-store", "must-revalidate"}},
			meth:        MethodPatchRaw,
			body_data:   []byte("this is byte"),
			want_status: 200,
		},
		{
			name:        "TEST STRING PATCH",
			url:         "http://httpbin.org/patch",
			params:      url.Values{"show_env": {"0"}},
			header:      http.Header{"Cache-Control": {"no-cache", "no-store", "must-revalidate"}},
			meth:        MethodPatchRaw,
			body_data:   "this is string",
			want_status: 200,
		},
		{
			name:        "TEST JSON PATCH",
			url:         "http://httpbin.org/patch",
			params:      url.Values{"show_env": {"0"}},
			header:      http.Header{"Cache-Control": {"no-cache", "no-store", "must-revalidate"}},
			meth:        MethodPatchJSON,
			body_data:   jsonex,
			want_status: 200,
		},
		{
			name:        "TEST XML PATCH",
			url:         "http://httpbin.org/patch",
			params:      url.Values{"show_env": {"0"}},
			header:      http.Header{"Cache-Control": {"no-cache", "no-store", "must-revalidate"}, "Content-Type": {"application/xml"}},
			meth:        MethodPatchXML,
			body_data:   xmlex,
			want_status: 200,
		},
		{
			name:        "TEST STANDARD DELETE",
			url:         "http://httpbin.org/delete",
			params:      url.Values{"show_env": {"0"}},
			header:      http.Header{"Cache-Control": {"no-cache", "no-store", "must-revalidate"}},
			meth:        MethodDelete,
			body_data:   nil,
			want_status: 200,
		},
		{
			name:        "TEST BYTE DELETE",
			url:         "http://httpbin.org/delete",
			params:      url.Values{"show_env": {"0"}},
			header:      http.Header{"Cache-Control": {"no-cache", "no-store", "must-revalidate"}},
			meth:        MethodDeleteRaw,
			body_data:   []byte("this is byte"),
			want_status: 200,
		},
		{
			name:        "TEST STRING DELETE",
			url:         "http://httpbin.org/delete",
			params:      url.Values{"show_env": {"0"}},
			header:      http.Header{"Cache-Control": {"no-cache", "no-store", "must-revalidate"}},
			meth:        MethodDeleteRaw,
			body_data:   "this is string",
			want_status: 200,
		},
		{
			name:        "TEST JSON DELETE",
			url:         "http://httpbin.org/delete",
			params:      url.Values{"show_env": {"0"}},
			header:      http.Header{"Cache-Control": {"no-cache", "no-store", "must-revalidate"}},
			meth:        MethodDeleteJSON,
			body_data:   jsonex,
			want_status: 200,
		},
		{
			name:        "TEST XML DELETE",
			url:         "http://httpbin.org/delete",
			params:      url.Values{"show_env": {"0"}},
			header:      http.Header{"Cache-Control": {"no-cache", "no-store", "must-revalidate"}, "Content-Type": {"application/xml"}},
			meth:        MethodDeleteXML,
			body_data:   xmlex,
			want_status: 200,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			header := tt.header
			params := tt.params
			httprequest := NewHTTPRequest()
			httprequest.Header = header
			httprequest.Method = tt.meth
			httprequest.URL = tt.url
			httprequest.Params = params
			httprequest.Body = tt.body_data

			response, _, _ := httprequest.Send()
			// fmt.Println("HTTP Response Status:", response.StatusCode, http.StatusText(response.StatusCode))
			// fmt.Println(string(body))
			assert.Equal(t, tt.want_status, response.StatusCode)
		})
	}
}
