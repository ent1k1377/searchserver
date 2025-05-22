package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type TestCase struct {
	name     string
	request  SearchRequest
	response SearchResponse
	err      error
}

func TestQwe(t *testing.T) {
	cases := []TestCase{
		{
			"deda",
			SearchRequest{
				Limit:      0,
				Offset:     0,
				Query:      "irure",
				OrderField: "",
				OrderBy:    0,
			},
			SearchResponse{},
			nil,
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	defer ts.Close()

	sc := SearchClient{
		AccessToken: "",
		URL:         ts.URL,
	}

	for _, item := range cases {
		t.Run(item.name, func(t *testing.T) {
			res, err := sc.FindUsers(item.request)

			if item.err != err {
				t.Errorf("Expected %v, got %v", item.err, err)
			}

			if item.response.NextPage != res.NextPage {
				t.Errorf("Expected %v, got %v", item.response.NextPage, res.NextPage)
			}
			// ...
		})
	}
}
