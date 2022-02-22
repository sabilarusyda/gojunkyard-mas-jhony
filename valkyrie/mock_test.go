package valkyrie

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

type httpFailClient struct{}

func (c *httpFailClient) Do(req *http.Request) (*http.Response, error) {
	return nil, http.ErrServerClosed
}

type http500Client struct{}

func (c *http500Client) Do(req *http.Request) (*http.Response, error) {
	const certs = `Internal Server Error`
	return &http.Response{
		Status:     http.StatusText(http.StatusInternalServerError),
		StatusCode: http.StatusInternalServerError,
		Body:       ioutil.NopCloser(bytes.NewReader([]byte(certs))),
	}, nil
}

type http200NilClient struct{}

func (c *http200NilClient) Do(req *http.Request) (*http.Response, error) {
	return &http.Response{
		Status:     http.StatusText(http.StatusOK),
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader(nil)),
	}, nil
}

type http200Client struct{}

func (c *http200Client) Do(req *http.Request) (*http.Response, error) {
	const hostmap = `{"data":[{"id":1,"type":"config","attributes":{"identifier":"supersoccertv_ads","name":"SSTV Ads","icon":"https://cdn01.supersoccer.tv/sstv/images/favicon.png","status":1,"host":["misty.supersoccer.tv/ads"],"modules":"[60,59,65,61,62,63,66,31,21,40,47,49,15,16,28,29,95,196,97,98,1,9,10,11,6,7,8,12,14,17,18,19,40,47,49,57,58,56,31,61,62,63,64,67,68,30,22,23,24,25,26,69,70,71,72,73,99,100,101,102,106,108,109,110,111,112,113,114,115,116,117,118,119]","created_at":"2018-10-19T09:08:17.000Z","updated_at":null,"deleted_at":null}},{"id":2,"type":"config","attributes":{"identifier":"molatv","name":"Mola Backoffice","icon":"https://cdn02.supersoccer.tv/Il3PgbLGR0GCvanrI9wU_misty-logo-32x32.png","status":1,"host":["misty.mola.id/back-office"],"modules":"[21,40,47,49,59,60,74,75,76,77,78,79,80,81,82,83,84,85,86,87,88,89,90,91,92,93,94]","created_at":"2018-10-05T15:37:20.000Z","updated_at":null,"deleted_at":null}},{"id":3,"type":"config","attributes":{"identifier":"supersoccertv_analytic","name":"SSTV Analytic","icon":"https://cdn01.supersoccer.tv/sstv/images/favicon.png","status":1,"host":["analytic.supersoccer.tv"],"modules":"[41,103,104]","created_at":"2018-10-15T16:49:03.000Z","updated_at":null,"deleted_at":null}}]}`
	return &http.Response{
		Status:     http.StatusText(http.StatusOK),
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader([]byte(hostmap))),
	}, nil
}
